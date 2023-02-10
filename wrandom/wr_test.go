// Copyright © 2021 Hedzr Yeh.

package wrandom_test

import (
	"testing"

	"github.com/hedzr/lb/lbapi"
	"github.com/hedzr/lb/random"
	"github.com/hedzr/lb/wrandom"
)

type exP string

func (s exP) String() string { return string(s) }

func withPeers(peers ...lbapi.Peer) lbapi.Opt {
	return func(balancer lbapi.Balancer) {
		for _, p := range peers {
			balancer.Add(p)
		}
	}
}

func TestWR1(t *testing.T) {
	peer1, peer2, peer3 := exP("172.16.0.7:3500"), exP("172.16.0.8:3500"), exP("172.16.0.9:3500")
	peer4, peer5 := exP("172.16.0.2:3500"), exP("172.16.0.3:3500")

	lb := wrandom.New(
		wrandom.WithWeightedBalancedPeers(
			wrandom.NewPeer(3, random.New, withPeers(peer1, peer2, peer3)),
			wrandom.NewPeer(2, random.New, withPeers(peer4, peer5)),
		),
	)

	sum := make(map[lbapi.Peer]int)

	for i := 0; i < 5000; i++ {
		p, _ := lb.Next(lbapi.DummyFactor)
		sum[p]++
	}

	// results
	for k, v := range sum {
		t.Logf("%v: %v", k, v)
	}
}
