// Copyright © 2021 Hedzr Yeh.

package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var port = 8103

func main() {
	target, err := url.Parse("http://ds1.service.local:8111")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("forwarding to -> %s://%s\n", target.Scheme, target.Host)

	proxy := httputil.NewSingleHostReverseProxy(target)

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		// https://stackoverflow.com/questions/38016477/reverse-proxy-does-not-work
		// https://forum.golangbridge.org/t/explain-how-reverse-proxy-work/6492/7
		// https://stackoverflow.com/questions/34745654/golang-reverseproxy-with-apache2-sni-hostname-error

		req.Host = req.URL.Host // if you remove this line the request will fail... I want to debug why.

		proxy.ServeHTTP(w, req)
	})

	fmt.Printf("Server started at port %v...\n", port)
	err = http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
	if err != nil {
		panic(err)
	}
}
