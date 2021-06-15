package wrandom

import (
	"github.com/hedzr/lb/lbapi"
	"github.com/hedzr/lb/wrr"
)

// New make a new instance of a weighted random balancer.
// But, you have had to pass a WithBalancedPeers opt to it.
//
// Example
//
//     peer1, peer2, peer3 := exP("172.16.0.7:3500"), exP("172.16.0.8:3500"), exP("172.16.0.9:3500")
//     peer4, peer5 := exP("172.16.0.2:3500"), exP("172.16.0.3:3500")
//     lb := wrandom.New(
//       wrandom.WithWeightedBalancedPeers(
//         wrandom.NewPeer(3, random.New, withPeers(peer1, peer2, peer3)),
//         wrandom.NewPeer(2, random.New, withPeers(peer4, peer5)),
//       ),
//     )
func New(opts ...lbapi.Opt) lbapi.Balancer { return wrr.New(opts...) }

func WithWeightedBalancedPeers(peers ...WBPeer) lbapi.Opt {
	return func(balancer lbapi.Balancer) {
		for _, b := range peers {
			if bp, ok := b.(lbapi.WeightedPeer); ok {
				balancer.Add(bp)
			}
		}
	}
}

func NewPeer(weight int, gen func(opts ...lbapi.Opt) lbapi.Balancer, opts ...lbapi.Opt) WBPeer {
	return &wpPeer{
		weight: weight,
		lb:     gen(opts...),
	}
}

// WBPeer is a weighted, balanced peer.
type WBPeer interface {
	lbapi.WeightedPeer
	lbapi.Balancer
}

type wpPeer struct {
	weight int
	lb     lbapi.Balancer
}

func (w *wpPeer) String() string { return "wpBeer" }

func (w *wpPeer) Weight() int { return w.weight }

func (w *wpPeer) Next(factor lbapi.Factor) (next lbapi.Peer, c lbapi.Constrainable) {
	return w.lb.Next(factor)
}
func (w *wpPeer) Count() int              { return w.lb.Count() }
func (w *wpPeer) Add(peers ...lbapi.Peer) { w.lb.Add(peers...) }
func (w *wpPeer) Remove(peer lbapi.Peer)  { w.lb.Remove(peer) }
func (w *wpPeer) Clear()                  { w.lb.Clear() }
