package rr

import (
	"github.com/hedzr/lb/lbapi"
	"sync"
	"sync/atomic"
)

// New make a new load-balancer instance with Round-Robin
func New(opts ...lbapi.Opt) lbapi.Balancer {
	return (&rrS{}).init(opts...)
}

type rrS struct {
	peers []lbapi.Peer
	count int64
	rw    sync.RWMutex
}

func (s *rrS) init(opts ...lbapi.Opt) *rrS {
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *rrS) Next(factor lbapi.Factor) (next lbapi.Peer, c lbapi.Constrainable) {
	next = s.miniNext()
	if fc, ok := factor.(lbapi.FactorComparable); ok {
		next, c, ok = fc.ConstrainedBy(next)
	} else if nested, ok := next.(lbapi.BalancerLite); ok {
		next, c = nested.Next(factor)
	}
	return
}

func (s *rrS) miniNext() (next lbapi.Peer) {
	ni := atomic.AddInt64(&s.count, 1)

	ni--

	s.rw.RLock()
	defer s.rw.RUnlock()
	ni %= int64(len(s.peers))
	next = s.peers[ni]
	return
}

func (s *rrS) Count() int {
	s.rw.RLock()
	defer s.rw.RUnlock()
	return len(s.peers)
}

func (s *rrS) Add(peers ...lbapi.Peer) {
	for _, p := range peers {
		s.AddOne(p)
	}
}

func (s *rrS) AddOne(peer lbapi.Peer) {
	if s.find(peer) {
		return
	}
	s.rw.Lock()
	defer s.rw.Unlock()
	s.peers = append(s.peers, peer)
}

func (s *rrS) find(peer lbapi.Peer) (found bool) {
	s.rw.RLock()
	defer s.rw.RUnlock()
	for _, p := range s.peers {
		if lbapi.DeepEqual(p, peer) {
			return true
		}
	}
	return
}

func (s *rrS) Remove(peer lbapi.Peer) {
	s.rw.Lock()
	defer s.rw.Unlock()
	for i, p := range s.peers {
		if lbapi.DeepEqual(p, peer) {
			s.peers = append(s.peers[0:i], s.peers[i+1:]...)
			return
		}
	}
}

func (s *rrS) Clear() {
	s.rw.Lock()
	defer s.rw.Unlock()
	s.peers = nil
}
