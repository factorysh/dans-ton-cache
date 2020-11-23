package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/factorysh/dans-ton-cache/cache"
)

func main() {
	rpURL, err := url.Parse(os.Getenv("BACKEND"))
	if err != nil {
		log.Fatal(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(rpURL)
	mux := http.NewServeMux()
	os.MkdirAll("/tmp/proxy", 0770)
	c, err := cache.New("/tmp/proxy", 10)
	if err != nil {
		log.Fatal(err)
	}
	mux.HandleFunc("/", c.Middleware(proxy.ServeHTTP))
	s := &http.Server{
		Addr:           os.Getenv("LISTEN"),
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}
