// Copyright © 2021 Hedzr Yeh.

package version

import (
	"sync"

	"github.com/hedzr/lb/lbapi"
)

// BackendsFactor interface
type BackendsFactor interface {
	lbapi.FactorComparable
	AddPeers(peers ...VersioningBackendFactor)
}

// NewBackendsFactor bundle a balancer generator and its opts into a BackendsFactor.
//
// A BackendsFactor holds a set of VersioningBackendFactor. These backends can be
// added with BackendsFactor.AddPeers(...).
//
// The bundled balancer will be used as a second-level balancer following up the parent.
//
// Example:
//
//	bf := version.NewBackendsFactor(rr.New)
//	bf.AddPeers(
//	    version.NewBackendFactor("1.1", "172.16.0.6:3500"),
//	    version.NewBackendFactor("1.3", "172.16.0.7:3500"),
//	    version.NewBackendFactor("2.0", "172.16.0.8:3500"),
//	    version.NewBackendFactor("3.13", "172.16.0.9:3500"),
//	)
//	//
//	var testConstraints = []lbapi.Constrainable{
//	    version.NewConstrainablePeer("<= 1.1.x", 2),
//	    version.NewConstrainablePeer("^1.2.x", 4),
//	    version.NewConstrainablePeer("^2.x", 11),
//	    version.NewConstrainablePeer("^3.x", 3),
//	}
//	lb := version.New(version.WithConstrainedPeers(testConstraints...))
//	factor := bf
//	peer, c := lb.Next(factor)
func NewBackendsFactor(gen func(opts ...lbapi.Opt) lbapi.Balancer, opts ...lbapi.Opt) BackendsFactor {
	return &backendsFactor{
		generator:   gen,
		opts:        opts,
		constraints: make(map[lbapi.Constrainable]lbapi.Balancer),
	}
}

type backendsFactor struct {
	backends    []VersioningBackendFactor
	constraints map[lbapi.Constrainable]lbapi.Balancer
	crw         sync.RWMutex
	generator   func(opts ...lbapi.Opt) lbapi.Balancer
	opts        []lbapi.Opt
}

func (fa *backendsFactor) AddPeers(peers ...VersioningBackendFactor) {
	fa.backends = append(fa.backends, peers...)
}

func (fa *backendsFactor) String() string { return "" }
func (fa *backendsFactor) Factor() string { return "" }
func (fa *backendsFactor) ConstrainedBy(constraints interface{}) (peer lbapi.Peer, c lbapi.Constrainable, satisfied bool) {
	if cc, ok := constraints.(lbapi.Constrainable); ok {
		var lb lbapi.Balancer

		// for this object cc, build a lb and associate with it
		fa.crw.RLock()
		if _, ok := fa.constraints[cc]; !ok {
			fa.crw.RUnlock()

			lb = fa.generator()
			fa.crw.Lock()
			fa.constraints[cc] = lb
			fa.crw.Unlock()
		} else {
			lb = fa.constraints[cc]
			fa.crw.RUnlock()
			lb.Clear()
		}

		// find all satisfied backends/peers and add them into lb
		for _, f := range fa.backends {
			if cc.Check(f) {
				satisfied = true
				lb.Add(f)
			}
		}

		// now, pick up the next peer of them
		peer, c = lb.Next(lbapi.DummyFactor)
		if c == nil {
			c = cc
		}
	}
	return
}
