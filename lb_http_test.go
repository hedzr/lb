// Copyright © 2021 Hedzr Yeh.

package lb

import (
	"context"
	"fmt"
	"github.com/hedzr/lb/lbapi"
	"github.com/hedzr/lb/wrr"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"testing"
)

func TestWebServers(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", helloHandler1)
	srv := &http.Server{
		Addr:    ":38080",
		Handler: mux,
	}

	mux = http.NewServeMux()
	mux.HandleFunc("/", helloHandler2)
	srv1 := &http.Server{
		Addr:    ":38081",
		Handler: mux,
	}

	var err error
	var exitCh = make(chan struct{})
	var wgStartup sync.WaitGroup
	wgStartup.Add(4)
	go func() {
		fmt.Println("Server started at port 38080")
		wgStartup.Done()
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			t.Errorf("listen: %+v\n", err)
		}
	}()
	go func() {
		fmt.Println("Server started at port 38081")
		wgStartup.Done()
		if err := srv1.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			t.Errorf("listen: %+v\n", err)
		}
	}()
	go func() {
		wgStartup.Done()
		<-exitCh
		if err = srv.Shutdown(context.Background()); err != nil {
			t.Errorf("server Shutdown Failed: %+v", err)
		}
		if err == http.ErrServerClosed {
			err = nil
		}
	}()
	go func() {
		wgStartup.Done()
		<-exitCh
		if err = srv1.Shutdown(context.Background()); err != nil {
			t.Errorf("server Shutdown Failed: %+v", err)
		}
		if err == http.ErrServerClosed {
			err = nil
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	sum := make(map[lbapi.Peer]int)
	res := make(map[string]int)
	go func() {
		defer wg.Done()
		wgStartup.Wait()
		t.Logf("requesting...")
		lb := wrr.New(wrr.WithPeersAndWeights(
			[]lbapi.Peer{
				&exP{"http://localhost:38080", 2},
				&exP{"http://localhost:38081", 3},
			},
			[]int{2, 3},
		))
		for i := 0; i < 300; i++ {
			p, _ := lb.Next(lbapi.DummyFactor)
			sum[p]++
			content, err := clientGet(t, p.String())
			if err == nil {
				res[string(content)]++
			} else {
				t.Errorf("client Get failed: %v", err)
			}
		}
		lb.Clear()
	}()

	wg.Wait()

	for k, v := range res {
		t.Logf("%v : %v", strings.TrimRight(k, "\r\n"), v)
	}
}

type exP struct {
	addr   string
	weight int
}

func (s *exP) String() string { return s.addr }
func (s *exP) Weight() int    { return s.weight }

func helloHandler1(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, there 38080\n")
}

func helloHandler2(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, there 38081\n")
}

func clientGet(t *testing.T, url string) (content []byte, err error) {
	var resp *http.Response
	resp, err = http.Get(url)
	if err != nil {
		t.Log("Get failed:", err)
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Log("statuscode:", resp.StatusCode)

	}

	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Log("Read failed:", err)
	}

	return
}
