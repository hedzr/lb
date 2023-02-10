// Copyright © 2021 Hedzr Yeh.

package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"

	"github.com/hedzr/lb/lbapi"
	"github.com/hedzr/lb/rr"
	"github.com/hedzr/log"
)

var port = 8103

type DebugTransport struct{}

func (DebugTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	b, err := httputil.DumpRequestOut(r, false)
	if err != nil {
		log.Errorf(" [proxy][api-gw] %v", err)
		return nil, err
	}
	log.Debugf(" [proxy][api-gw] %v", string(b))
	return http.DefaultTransport.RoundTrip(r)
}

type ProxyPeer struct {
	*httputil.ReverseProxy
	url    string
	weight int
}

func (p ProxyPeer) String() string { return p.String() }
func (p ProxyPeer) Weight() int    { return p.weight }

func main() {
	log.SetLevel(log.DebugLevel)

	var ports []int
	for i := 1; i < len(os.Args); i++ {
		str := os.Args[i]
		i, err := strconv.Atoi(str)
		if err == nil {
			ports = append(ports, i)
		}
	}
	if len(ports) == 0 {
		ports = []int{8111, 8112}
	}

	var b = rr.New()
	for _, p := range ports {
		urlTarget := fmt.Sprintf("%s://ds1.service.local:%v", "http", p)
		target, err := url.Parse(urlTarget)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("forwarding to -> %s\n", target)
		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.Transport = DebugTransport{}
		b.Add(&ProxyPeer{proxy, urlTarget, 1})
	}

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		req.Host = req.URL.Host
		peer, _ := b.Next(lbapi.DummyFactor)
		peer.(*ProxyPeer).ServeHTTP(w, req)
	})

	fmt.Printf("Server started at port %v...\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
	if err != nil {
		panic(err)
	}
}
