// Copyright © 2021 Hedzr Yeh.

package version

import (
	"github.com/hedzr/lb/lbapi"
	"github.com/hedzr/lb/wrr"
)

// New make a new versioning wrr balancer.
//
// Example:
//
//    bf := version.NewBackendsFactor(rr.New)
//    bf.AddPeers(
//        version.NewBackendFactor("1.1", "172.16.0.6:3500"),
//        version.NewBackendFactor("1.3", "172.16.0.7:3500"),
//        version.NewBackendFactor("2.0", "172.16.0.8:3500"),
//        version.NewBackendFactor("3.13", "172.16.0.9:3500"),
//    )
//    //
//    var testConstraints = []lbapi.Constrainable{
//        version.NewConstrainablePeer("<= 1.1.x", 2),
//        version.NewConstrainablePeer("^1.2.x", 4),
//        version.NewConstrainablePeer("^2.x", 11),
//        version.NewConstrainablePeer("^3.x", 3),
//    }
//    lb := version.New(version.WithConstrainedPeers(testConstraints...))
//    factor := bf
//    peer, c := lb.Next(factor)
//
//
func New(opts ...lbapi.Opt) lbapi.Balancer { return wrr.New(opts...) }

// WithConstrainedPeers fills a set of lbapi.Constrainable as peers.
//
// Example:
//
//    bf := version.NewBackendsFactor(rr.New)
//    bf.AddPeers(
//        version.NewBackendFactor("1.1", "172.16.0.6:3500"),
//        version.NewBackendFactor("1.3", "172.16.0.7:3500"),
//        version.NewBackendFactor("2.0", "172.16.0.8:3500"),
//        version.NewBackendFactor("3.13", "172.16.0.9:3500"),
//    )
//    //
//    var testConstraints = []lbapi.Constrainable{
//        version.NewConstrainablePeer("<= 1.1.x", 2),
//        version.NewConstrainablePeer("^1.2.x", 4),
//        version.NewConstrainablePeer("^2.x", 11),
//        version.NewConstrainablePeer("^3.x", 3),
//    }
//    lb := version.New(version.WithConstrainedPeers(testConstraints...))
//    factor := bf
//    peer, c := lb.Next(factor)
//
func WithConstrainedPeers(cs ...lbapi.Constrainable) lbapi.Opt {
	return func(balancer lbapi.Balancer) {
		for _, vp := range cs {
			balancer.Add(vp)
		}
	}
}
