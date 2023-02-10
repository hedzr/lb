// Copyright © 2021 Hedzr Yeh.

// Package lb provides a flexible load balancer with pluggable strategies.
package lb

import (
	"log"
	"sync"

	"github.com/hedzr/lb/hash"
	"github.com/hedzr/lb/lbapi"
	"github.com/hedzr/lb/random"
	"github.com/hedzr/lb/rr"
	"github.com/hedzr/lb/version"
	"github.com/hedzr/lb/wrandom"
	"github.com/hedzr/lb/wrr"
)

// New make a new instance of a balancer.
//
//	l := lb.New(lb.WeightedRoundRobin, lb.WithPeers(some-peers-here...))
//	fmt.Println(l.Next(lbapi.DummyFactor)
//
// check out the real example in test codes.
func New(algorithm string, opts ...lbapi.Opt) lbapi.Balancer {
	kbs.RLock()
	defer kbs.RUnlock()
	if g, ok := knownBalancers[algorithm]; ok {
		return g(opts...)
	}
	log.Fatalf("unknown/unregistered balancer and generator: %q", algorithm)
	return nil // unreachable
}

// WithPeers adds the initial peers.
func WithPeers(peers ...lbapi.Peer) lbapi.Opt {
	return func(balancer lbapi.Balancer) {
		for _, p := range peers {
			balancer.Add(p)
		}
	}
}

// Register assign a (algorithm, generator) pair.
func Register(algorithm string, generator func(opts ...lbapi.Opt) lbapi.Balancer) {
	kbs.Lock()
	defer kbs.Unlock()
	knownBalancers[algorithm] = generator
}

// Unregister revoke a (algorithm, generator) pair.
func Unregister(algorithm string) {
	kbs.Lock()
	defer kbs.Unlock()
	delete(knownBalancers, algorithm)
}

const (
	// Random algorithm
	Random = "random"
	// RoundRobin algorithm
	RoundRobin = "round-robin"
	// WeightedRoundRobin algorithm
	WeightedRoundRobin = "weighted-round-robin"
	// ConsistentHash algorithm
	ConsistentHash = "consistent-hash"
	// WeightedRandom algorithm
	WeightedRandom = "weighted-random"
	// VersioningWRR algorithm
	VersioningWRR = "versioning-wrr"
)

func init() {
	kbs.Lock()
	defer kbs.Unlock()

	knownBalancers = make(map[string]func(opts ...lbapi.Opt) lbapi.Balancer)

	knownBalancers[Random] = random.New
	knownBalancers[RoundRobin] = rr.New
	knownBalancers[WeightedRoundRobin] = wrr.New
	knownBalancers[ConsistentHash] = hash.New

	knownBalancers[WeightedRandom] = wrandom.New

	knownBalancers[VersioningWRR] = version.New
}

var knownBalancers map[string]func(opts ...lbapi.Opt) lbapi.Balancer
var kbs sync.RWMutex

// need await for go 2 generic
