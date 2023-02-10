// Copyright © 2021 Hedzr Yeh.

package wrr_test

import (
	"sync"
	"testing"

	"github.com/hedzr/lb/lbapi"
	"github.com/hedzr/lb/wrr"
)

type exP struct {
	addr   string
	weight int
}

func (s *exP) String() string { return s.addr }
func (s *exP) Weight() int    { return s.weight }

func TestWRR1(t *testing.T) {
	lb := wrr.New()
	lb.Add(&exP{"172.16.0.7:3500", 5}, &exP{"172.16.0.8:3500", 3}, &exP{"172.16.0.9:3500", 2})

	sum := make(map[lbapi.Peer]int)

	for i := 0; i < 300; i++ {
		p, _ := lb.Next(lbapi.DummyFactor)
		sum[p]++
	}

	// results
	total := 0
	for _, v := range sum {
		total += v
	}
	for k, v := range sum {
		t.Logf("%v: weight = %v, %v/%0.2f%%", k, k.(lbapi.WeightedPeer).Weight(), v, (float32(v)/float32(total))*100.0)
	}
}

func TestWRR2(t *testing.T) {
	lb := wrr.New(wrr.WithWeightedPeers(&exP{"172.16.0.7:3500", 5}))
	sum := make(map[lbapi.Peer]int)
	for i := 0; i < 300; i++ {
		p, _ := lb.Next(lbapi.DummyFactor)
		sum[p]++
	}
}

func TestWRR3(t *testing.T) {
	lb := wrr.New(wrr.WithPeersAndWeights([]lbapi.Peer{&exP{"abd", 2}, &exP{"zyx", 3}}, []int{2, 3}))
	sum := make(map[lbapi.Peer]int)
	for i := 0; i < 300; i++ {
		p, _ := lb.Next(lbapi.DummyFactor)
		sum[p]++
	}
	lb.Clear()
}

func TestWRR_AddRemove(t *testing.T) {
	lb := wrr.New()
	lb.Add(&exP{"172.16.0.7:3500", 5}, &exP{"172.16.0.8:3500", 3}, &exP{"172.16.0.9:3500", 2})

	lb.Add(&exP{"172.16.0.8:3500", 3})
	if lb.Count() != 3 {
		t.Fatal("wrong Add: the dup peer should be ignore")
	}

	lb.Remove(&exP{"172.16.0.8:3500", 3})
	if lb.Count() != 2 {
		t.Fatalf("wrong Remove: not removed? count = %v", lb.Count())
	}
}

func adder(key lbapi.Peer, sum map[lbapi.Peer]int, rw *sync.RWMutex) {
	rw.Lock()
	defer rw.Unlock()
	sum[key]++
}

func TestWRR_M1(t *testing.T) {
	lb := wrr.New()
	lb.Add(&exP{"172.16.0.7:3500", 5}, &exP{"172.16.0.8:3500", 3}, &exP{"172.16.0.9:3500", 2})

	var wg sync.WaitGroup
	var rw sync.RWMutex
	sum := make(map[lbapi.Peer]int)

	const threads = 8
	wg.Add(threads)
	for x := 0; x < threads; x++ {
		go func(xi int) {
			defer wg.Done()
			for i := 0; i < 600; i++ {
				p, _ := lb.Next(lbapi.DummyFactor)
				adder(p, sum, &rw)
			}
		}(x)
	}

	wg.Wait()

	// results
	total := 0
	for _, v := range sum {
		total += v
	}
	for k, v := range sum {
		t.Logf("%v: weight = %v, %v/%0.2f%%", k, k.(lbapi.WeightedPeer).Weight(), v, (float32(v)/float32(total))*100.0)
	}
}
