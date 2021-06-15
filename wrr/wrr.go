package wrr

import (
	"github.com/hedzr/lb/lbapi"
	"sync"
)

// New make a new load-balancer instance with Weighted Round-Robin
func New(opts ...lbapi.Opt) lbapi.Balancer {
	return (&wrrS{
		m: make(map[lbapi.Peer]*weightS),
	}).init(opts...)
}

// WithPeersAndWeights allows passing a simple peer array and
// corresponding weight array separately.
func WithPeersAndWeights(peers []lbapi.Peer, weights []int) lbapi.Opt {
	return func(balancer lbapi.Balancer) {
		if wrr, ok := balancer.(*wrrS); ok {
			wrr.addWeights(peers, weights)
		}
	}
}

// WithWeightedPeers allows passing a weighted-peers array.
func WithWeightedPeers(peers ...lbapi.WeightedPeer) lbapi.Opt {
	return func(balancer lbapi.Balancer) {
		if wrr, ok := balancer.(*wrrS); ok {
			wrr.addPeers(peers)
		}
	}
}

type wrrS struct {
	peers []lbapi.Peer
	m     map[lbapi.Peer]*weightS
	prw   sync.RWMutex // for peers
	mrw   sync.RWMutex // for m
}

type weightS struct {
	weight    int
	effective int
	current   int
}

func (s *wrrS) init(opts ...lbapi.Opt) *wrrS {
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Next implements a smooth weighted round-robin lb with algorithm coming from nginx:
// https://github.com/nginx/nginx/commit/52327e0627f49dbda1e8db695e63a4b0af4448b1
func (s *wrrS) Next(factor lbapi.Factor) (best lbapi.Peer, c lbapi.Constrainable) {
	if best = s.miniNext(); best != nil {
		if fc, ok := factor.(lbapi.FactorComparable); ok {
			best, c, _ = fc.ConstrainedBy(best)
		} else if nested, ok := best.(lbapi.BalancerLite); ok {
			best, c = nested.Next(factor)
		}
	}
	return
}

func (s *wrrS) miniNext() (best lbapi.Peer) {
	total := 0

	s.prw.RLock()
	defer s.prw.RUnlock()

	for _, node := range s.peers {
		if node == nil {
			continue
		}

		total += s.mUpdate(node, 0, false)
		if s.mTest(best, node) {
			best = node
		}
	}

	if best != nil {
		s.mUpdate(best, -total, true)
	}
	return
}

func (s *wrrS) mUpdate(node lbapi.Peer, delta int, success bool) (total int) {
	s.mrw.Lock()
	defer s.mrw.Unlock()
	if delta == 0 {
		delta = s.m[node].effective
	}
	s.m[node].current += delta
	//if success {
	//	s.m[node].effective++
	//}
	return s.m[node].effective
}

func (s *wrrS) mTest(best, node lbapi.Peer) bool {
	s.mrw.RLock()
	defer s.mrw.RUnlock()
	return best == nil || s.m[node].current > s.m[best].current
}

func (s *wrrS) addPeers(peers []lbapi.WeightedPeer) {
	s.prw.RLock()
	defer s.prw.RUnlock()
	for _, p := range peers {
		s.peers = append(s.peers, p)
	}

	s.mrw.Lock()
	defer s.mrw.Unlock()

	for _, p := range peers {
		var w = p.Weight()
		s.m[p] = &weightS{current: 0, effective: w, weight: w}
	}
}

func (s *wrrS) addWeights(peers []lbapi.Peer, weights []int) {
	s.prw.RLock()
	defer s.prw.RUnlock()
	s.peers = peers

	s.mrw.Lock()
	defer s.mrw.Unlock()

	for i, p := range peers {
		var w int
		if wp, ok := p.(lbapi.WeightedPeer); ok {
			w = wp.Weight()
		} else {
			w = weights[i]
		}
		s.m[p] = &weightS{current: 0, effective: w, weight: w}
	}
}

func (s *wrrS) SetNodeWeight(node lbapi.WeightedPeer, newWeight int) {
	s.mAdd(node, newWeight)
}

func (s *wrrS) mAdd(node lbapi.Peer, weight int) {
	s.mrw.Lock()
	defer s.mrw.Unlock()

	if wp, ok := node.(lbapi.WeightedPeer); ok {
		weight = wp.Weight()
	}

	if v, ok := s.m[node]; ok {
		v.effective = weight
		if v.weight == 0 {
			v.weight = weight
		}
	} else {
		s.m[node] = &weightS{current: 0, effective: weight, weight: weight}
	}
}

func (s *wrrS) Count() int {
	s.prw.RLock()
	defer s.prw.RUnlock()
	return len(s.peers)
}

func (s *wrrS) Add(peers ...lbapi.Peer) {
	for _, p := range peers {
		s.AddOne(p)
	}
}

func (s *wrrS) AddOne(peer lbapi.Peer) {
	if s.find(peer) {
		return
	}

	s.prw.Lock()
	defer s.prw.Unlock()
	s.peers = append(s.peers, peer)
	if wp, ok := peer.(lbapi.WeightedPeer); ok {
		s.SetNodeWeight(wp, wp.Weight())
	} else if _, ok := s.m[peer]; !ok {
		s.m[peer] = &weightS{
			weight:    0,
			effective: 0,
			current:   0,
		}
	}
}

func (s *wrrS) find(peer lbapi.Peer) (found bool) {
	s.prw.RLock()
	defer s.prw.RUnlock()
	for _, p := range s.peers {
		if lbapi.DeepEqual(p, peer) {
			return true
		}
	}
	return
}

func (s *wrrS) Remove(peer lbapi.Peer) {
	s.prw.Lock()
	defer s.prw.Unlock()
	for i, p := range s.peers {
		if lbapi.DeepEqual(p, peer) {
			s.peers = append(s.peers[0:i], s.peers[i+1:]...)
			return
		}
	}
}

func (s *wrrS) Clear() {
	s.prw.Lock()
	defer s.prw.Unlock()
	s.peers = nil

	s.mrw.Lock()
	defer s.mrw.Unlock()
	s.m = make(map[lbapi.Peer]*weightS)
}
