package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func main() {
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   "127.0.0.1:8000",
		Path:   "/",
	})

	addr := "127.0.0.1:8080"
	fmt.Println("Listening to", addr)
	if err := http.ListenAndServe(addr, proxy); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
