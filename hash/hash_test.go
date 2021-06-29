// Copyright © 2021 Hedzr Yeh.

package hash_test

import (
	"github.com/hedzr/lb/hash"
	"github.com/hedzr/lb/lbapi"
	"hash/crc32"
	"strings"
	"sync"
	"testing"
)

type exP string

func (s exP) String() string { return string(s) }

var factors = []lbapi.FactorString{
	"https://abc.local/user/profile",
	"https://abc.local/admin/",
	"https://abc.local/shop/item/1",
	"https://abc.local/post/35719",
}

func TestHash1(t *testing.T) {
	lb := hash.New()
	lb.Add(exP("172.16.0.7:3500"), exP("172.16.0.8:3500"), exP("172.16.0.9:3500"))

	sum := make(map[lbapi.Peer]int)
	hits := make(map[lbapi.Peer]map[lbapi.Factor]bool)

	for i := 0; i < 300; i++ {
		factor := factors[i%len(factors)]
		peer, _ := lb.Next(factor)

		sum[peer]++
		if ps, ok := hits[peer]; ok {
			if _, ok := ps[factor]; !ok {
				ps[factor] = true
			}
		} else {
			hits[peer] = make(map[lbapi.Factor]bool)
			hits[peer][factor] = true
		}
	}

	// results
	total := 0
	for _, v := range sum {
		total += v
	}
	for p, v := range sum {
		var keys []string
		for fs := range hits[p] {
			if kk, ok := fs.(interface{ String() string }); ok {
				keys = append(keys, kk.String())
			} else {
				keys = append(keys, fs.Factor())
			}
		}
		t.Logf("%v: %v, [%v]", p, v, strings.Join(keys, ","))
	}

	lb.Clear()
}

func adder(key lbapi.Peer, c lbapi.Constrainable, sum map[lbapi.Peer]int, rw *sync.RWMutex) {
	rw.Lock()
	defer rw.Unlock()
	sum[key]++
}

func TestHash_AddRemove(t *testing.T) {
	lb := hash.New()
	lb.Add(exP("172.16.0.7:3500"), exP("172.16.0.8:3500"), exP("172.16.0.9:3500"))

	lb.Add(exP("172.16.0.8:3500"))
	if lb.Count() != 3 {
		t.Fatal("wrong Add: the dup peer should be ignore")
	}
	lb.Remove(exP("172.16.0.8:3500"))
	if lb.Count() != 2 {
		t.Fatalf("wrong Remove: not removed? count = %v", lb.Count())
	}
}

func TestHash_M1(t *testing.T) {
	lb := hash.New()
	lb.Add(exP("172.16.0.7:3500"), exP("172.16.0.8:3500"), exP("172.16.0.9:3500"))

	var wg sync.WaitGroup
	var rw sync.RWMutex
	sum := make(map[lbapi.Peer]int)

	const threads = 8
	wg.Add(threads)
	for x := 0; x < threads; x++ {
		go func(xi int) {
			defer wg.Done()
			for i := 0; i < 600; i++ {
				p, c := lb.Next(factors[i%3])
				adder(p, c, sum, &rw)
			}
		}(x)
	}

	wg.Wait()

	// results
	for k, v := range sum {
		t.Logf("%v: %v", k, v)
	}
}

func TestHash2(t *testing.T) {
	lb := hash.New(
		hash.WithHashFunc(crc32.ChecksumIEEE),
		hash.WithReplica(16),
	)
	lb.Add(exP("172.16.0.7:3500"), exP("172.16.0.8:3500"), exP("172.16.0.9:3500"))

	sum := make(map[lbapi.Peer]int)
	hits := make(map[lbapi.Peer]map[lbapi.Factor]bool)

	for i := 0; i < 300; i++ {
		factor := factors[i%len(factors)]
		peer, _ := lb.Next(factor)

		sum[peer]++
		if ps, ok := hits[peer]; ok {
			if _, ok := ps[factor]; !ok {
				ps[factor] = true
			}
		} else {
			hits[peer] = make(map[lbapi.Factor]bool)
			hits[peer][factor] = true
		}
	}

	lb.Clear()
}
