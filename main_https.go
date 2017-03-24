package main

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"

	"cloud.google.com/go/storage"

	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/net/context"
)

func https() *http.Server {
	ctx := context.Background()
	gcs, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Could not make a storage client: %v", err)
	}
	cache := &GCSCache{
		Bucket: gcs.Bucket("cbro-scratch"),
		Prefix: "autocert-cache/autocert/",
	}

	m := &autocert.Manager{
		Cache:  cache,
		Prompt: autocert.AcceptTOS,
		Email:  "cbro@golang.org",
	}

	return &http.Server{
		TLSConfig: &tls.Config{
			GetCertificate: m.GetCertificate,
		},
	}
}

type GCSCache struct {
	Bucket *storage.BucketHandle
	Prefix string
}

func (c *GCSCache) Get(ctx context.Context, key string) (data []byte, err error) {
	r, err := c.Bucket.Object(c.Prefix + key).NewReader(ctx)
	if err == storage.ErrObjectNotExist {
		return nil, autocert.ErrCacheMiss
	}
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return ioutil.ReadAll(r)
}

func (c *GCSCache) Put(ctx context.Context, key string, data []byte) error {
	w := c.Bucket.Object(c.Prefix + key).NewWriter(ctx)
	if _, err := w.Write(data); err != nil {
		return err
	}
	return w.Close()
}

func (c *GCSCache) Delete(ctx context.Context, key string) error {
	return c.Bucket.Object(c.Prefix + key).Delete(ctx)
}
