package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"time"
)

func init() {
	// Simulate a slow app start-up.
	time.Sleep(5 * time.Second)
}

var serverStart = time.Now()

func main() {
	http.HandleFunc("/", handle)

	const addr = ":8080"
	srv := http.Server{}
	srv.Addr = addr
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

/*
Move to main() to turn on HTTPS.

	go func() {
		const addr = ":8083"
		log.Printf("Listening on https://%s", addr)
		srv := https()
		srv.Addr = addr
		log.Fatal(srv.ListenAndServeTLS("", ""))
	}()

*/
