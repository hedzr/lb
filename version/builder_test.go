// Copyright © 2021 Hedzr Yeh.

package version_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hedzr/lb/lbapi"
	"github.com/hedzr/lb/rr"
	"github.com/hedzr/lb/version"
	"github.com/hedzr/lb/wrr"
)

func initFactors() version.BackendsFactor {
	fa := version.NewBackendsFactor(rr.New)
	fa.AddPeers(
		version.NewBackendFactor("1.1", "172.16.0.6:3500"),
		version.NewBackendFactor("1.3", "172.16.0.7:3500"),
		version.NewBackendFactor("2.0", "172.16.0.8:3500"),
		version.NewBackendFactor("3.13", "172.16.0.9:3500"),
	)
	return fa
}

//func findC(peer lbapi.Peer) *exPeer {
//	for _, expr := range testConstraints {
//		if ex, ok := expr.(*exPeer); ok {
//			if ex.Check(peer) {
//				return ex
//			}
//		}
//	}
//	return nil
//}

func TestVersionWRRBuilder(t *testing.T) {
	b := version.NewBuilder()
	for _, vp := range testConstraints {
		b.AddPeer(vp.String(), vp.(lbapi.Weighted).Weight())
	}
	b.AddPeers(testConstraints...)
	lb := b.Build()

	factor := initFactors()
	for i := 0; i < 300; i++ {
		_, _ = lb.Next(factor)
	}
}

func TestVersionWRR1(t *testing.T) {
	factor := initFactors()

	lb := wrr.New()
	for _, vp := range testConstraints {
		lb.Add(vp)
	}

	sum := make(map[lbapi.Peer]int)
	hits := make(map[lbapi.Peer]map[lbapi.Constrainable]bool)

	for i := 0; i < 300; i++ {
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
