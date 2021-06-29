// Copyright © 2021 Hedzr Yeh.

package main

import (
	"fmt"
	"github.com/hedzr/lb/lbapi"
	mrand "math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type randomS struct {
	peers []lbapi.Peer
	count int64
}

func (s *randomS) Next(factor lbapi.Factor) (next lbapi.Peer, c lbapi.Constrainable) {
	l := int64(len(s.peers))
	ni := atomic.AddInt64(&s.count, inRange(0, l)) % l
	next = s.peers[ni]
	return
}

func main() {
	lb := &randomS{
		peers: []lbapi.Peer{
			exP("172.16.0.7:3500"), exP("172.16.0.8:3500"), exP("172.16.0.9:3500"),
		},
		count: 0,
	}

	sum := make(map[lbapi.Peer]int)
	for i := 0; i < 300; i++ {
		p, _ := lb.Next(lbapi.DummyFactor)
		sum[p]++
	}

	for k, v := range sum {
		fmt.Printf("%v: %v\n", k, v)
	}
}

var seededRand = mrand.New(mrand.NewSource(time.Now().UnixNano()))
var seedmu sync.Mutex

func inRange(min, max int64) int64 {
	seedmu.Lock()
	defer seedmu.Unlock()
	return seededRand.Int63n(max-min) + min
}

type exP string

func (s exP) String() string { return string(s) }
