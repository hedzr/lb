package random_test

import (
	"github.com/hedzr/lb/lbapi"
	"github.com/hedzr/lb/random"
	"sync"
	"testing"
)

type exP string

func (s exP) String() string { return string(s) }

func TestRandom1(t *testing.T) {
	lb := random.New()
	lb.Add(exP("172.16.0.7:3500"), exP("172.16.0.8:3500"), exP("172.16.0.9:3500"))

	sum := make(map[lbapi.Peer]int)

	for i := 0; i < 300; i++ {
		p, _ := lb.Next(lbapi.DummyFactor)
		sum[p]++
	}

	// results
	for k, v := range sum {
		t.Logf("%v: %v", k, v)
	}
}

func adder(key lbapi.Peer, sum map[lbapi.Peer]int, rw *sync.RWMutex) {
	rw.Lock()
	defer rw.Unlock()
	sum[key]++
}

func TestRandom_AddRemove(t *testing.T) {
	lb := random.New()
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

func TestRandom_M1(t *testing.T) {
	lb := random.New()
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
				p, _ := lb.Next(lbapi.DummyFactor)
				adder(p, sum, &rw)
			}
		}(x)
	}

	wg.Wait()

	// results
	for k, v := range sum {
		t.Logf("%v: %v", k, v)
	}

	lb.Clear()
}
