package random

import (
	"github.com/hedzr/lb/lbapi"
	mrand "math/rand"
	"sync"
	"sync/atomic"
	"time"
)

var seededRand = mrand.New(mrand.NewSource(time.Now().UnixNano()))
var seedmu sync.Mutex

func inRange(min, max int64) int64 {
	seedmu.Lock()
	defer seedmu.Unlock()
	return seededRand.Int63n(max-min) + min
}

// New make a new load-balancer instance with Round-Robin
func New(opts ...lbapi.Opt) lbapi.Balancer {
	return (&randomS{}).init(opts...)
}

type randomS struct {
	peers []lbapi.Peer
	count int64
	rw    sync.RWMutex
}

func (s *randomS) init(opts ...lbapi.Opt) *randomS {
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *randomS) String() string { return "random" }

func (s *randomS) Next(factor lbapi.Factor) (next lbapi.Peer, c lbapi.Constrainable) {
	next = s.miniNext()
	if fc, ok := factor.(lbapi.FactorComparable); ok {
		next, c, ok = fc.ConstrainedBy(next)
	} else if nested, ok := next.(lbapi.BalancerLite); ok {
		next, c = nested.Next(factor)
	}
	return
}

func (s *randomS) miniNext() (next lbapi.Peer) {
	s.rw.RLock()
	defer s.rw.RUnlock()

	l := int64(len(s.peers))
	ni := atomic.AddInt64(&s.count, inRange(0, l)) % l
	next = s.peers[ni]
	return
}

func (s *randomS) Count() int {
	s.rw.RLock()
	defer s.rw.RUnlock()
	return len(s.peers)
}

func (s *randomS) Add(peers ...lbapi.Peer) {
	for _, p := range peers {
		s.AddOne(p)
	}
}

func (s *randomS) AddOne(peer lbapi.Peer) {
	if s.find(peer) {
		return
	}

	s.rw.Lock()
	defer s.rw.Unlock()
	s.peers = append(s.peers, peer)
}

func (s *randomS) find(peer lbapi.Peer) (found bool) {
	s.rw.RLock()
	defer s.rw.RUnlock()
	for _, p := range s.peers {
		if lbapi.DeepEqual(p, peer) {
			return true
		}
	}
	return
}

func (s *randomS) Remove(peer lbapi.Peer) {
	s.rw.Lock()
	defer s.rw.Unlock()
	for i, p := range s.peers {
		if lbapi.DeepEqual(p, peer) {
			s.peers = append(s.peers[0:i], s.peers[i+1:]...)
			return
		}
	}
}

func (s *randomS) Clear() {
	s.rw.Lock()
	defer s.rw.Unlock()
	s.peers = nil
}
