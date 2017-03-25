package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"syscall"
	"time"
)

func init() {
	// Simulate a slow app start-up.
	time.Sleep(6 * time.Second)
}

var serverStart = time.Now()

func main() {
	http.HandleFunc("/", handle)

	go func() {
		const addr = ":8083"
		srv := https()
		srv.Addr = addr
		shutdownOnSignal(srv, syscall.SIGTERM)
		log.Printf("Listening on https://%s", addr)
		log.Fatal(srv.ListenAndServeTLS("", ""))
	}()

	const addr = ":8080"
	srv := &http.Server{}
	srv.Addr = addr
	shutdownOnSignal(srv, syscall.SIGTERM)
	log.Printf("Listening on http://%s", addr)
	log.Fatal(srv.ListenAndServe())
}

func handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Uptime: %s\n", time.Now().Sub(serverStart))

	env := os.Environ()
	sort.Strings(env)
	for _, e := range env {
		fmt.Fprintln(w, e)
	}
}

func shutdownOnSignal(srv *http.Server, sig ...os.Signal) {
	sigc := make(chan os.Signal)
	signal.Notify(sigc, sig...)
	go func() {
		<-sigc
		log.Printf("Shutting down %v", srv.Addr)
		ctx := context.Background()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Could not shut down %v: %v", srv.Addr, err)
		} else {
			log.Printf("Shut down %v cleanly", srv.Addr)
		}
	}()
}

/*
Move to main() to turn on HTTPS.

	go func() {
		const addr = ":8083"
		srv := https()
		srv.Addr = addr
		log.Printf("Listening on https://%s", addr)
		log.Fatal(srv.ListenAndServeTLS("", ""))
	}()

*/

/*
Add to each HTTP server to shut down cleanly on SIGTERM.

		shutdownOnSignal(srv, syscall.SIGTERM)

*/
