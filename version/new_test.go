// Copyright © 2021 Hedzr Yeh.

package version_test

import (
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/hedzr/lb/lbapi"
	"github.com/hedzr/lb/version"
)

var testConstraints = []lbapi.Constrainable{
	version.NewConstrainablePeer("<= 1.1.x", 2),
	version.NewConstrainablePeer("^1.2.x", 4),
	version.NewConstrainablePeer("^2.x", 11),
	version.NewConstrainablePeer("^3.x", 3),
}

func TestVersionWRR2(t *testing.T) {
	lb := version.New(version.WithConstrainedPeers(testConstraints...))

	sum := make(map[lbapi.Peer]int)
	hits := make(map[lbapi.Peer]map[lbapi.Constrainable]bool)

	factor := initFactors()

	for i := 0; i < 500; i++ {
		peer, c := lb.Next(factor)

		sum[peer]++
		if ps, ok := hits[peer]; ok {
			if _, ok := ps[c]; !ok {
				ps[c] = true
			}
		} else {
			hits[peer] = make(map[lbapi.Constrainable]bool)
			hits[peer][c] = true
		}
	}

	// results
	total := 0
	for _, v := range sum {
		total += v
	}
	for peer, v := range sum {
		var keys []string
		var w int
		for c := range hits[peer] {
			if kk, ok := c.(fmt.Stringer); ok {
				keys = append(keys, kk.String())
			}
			if ww, ok := c.(lbapi.Weighted); ok {
				w = ww.Weight()
			}
		}
		// ex := findC(peer)
		t.Logf("%v: %v/%0.2f%%/w:%v. [%v => weight: %v]",
			peer, v, (float32(v)/float32(total))*100.0, w,
			strings.Join(keys, ","), w)
	}
}

func TestVersionWRR2_M1(t *testing.T) {
	lb := version.New(version.WithConstrainedPeers(testConstraints...))

	var wg sync.WaitGroup
	var rw sync.RWMutex
	sum := make(map[lbapi.Peer]int)
	hits := make(map[lbapi.Peer]map[lbapi.Constrainable]bool)

	const threads = 8
	wg.Add(threads)
	for x := 0; x < threads; x++ {
		go func(xi int) {
			defer wg.Done()
			for i := 0; i < 600; i++ {
				p, c := lb.Next(lbapi.DummyFactor)
				adder(p, sum, &rw, hits, c)
			}
		}(x)
	}

	wg.Wait()

	// results
	total := 0
	for _, v := range sum {
		total += v
	}
	//for k, v := range sum {
	//	t.Logf("%v: weight = %v, %v/%0.2f%%", k, k.(lbapi.WeightedPeer).Weight(), v, (float32(v)/float32(total))*100.0)
	//}
	for peer, v := range sum {
		var keys []string
		var w int
		for c := range hits[peer] {
			if kk, ok := c.(fmt.Stringer); ok {
				keys = append(keys, kk.String())
			}
			if ww, ok := c.(lbapi.Weighted); ok {
				w = ww.Weight()
			}
		}
		// ex := findC(peer)
		t.Logf("%v: %v/%0.2f%%/w:%v. [%v => weight: %v]",
			peer, v, (float32(v)/float32(total))*100.0, w,
			strings.Join(keys, ","), w)
	}
}

func adder(key lbapi.Peer, sum map[lbapi.Peer]int, rw *sync.RWMutex, hits map[lbapi.Peer]map[lbapi.Constrainable]bool, c lbapi.Constrainable) {
	rw.Lock()
	defer rw.Unlock()

	sum[key]++

	if ps, ok := hits[key]; ok {
		if _, ok := ps[c]; !ok {
			ps[c] = true
		}
	} else {
		hits[key] = make(map[lbapi.Constrainable]bool)
		hits[key][c] = true
	}
}
