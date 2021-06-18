package hash

import (
	"fmt"
	"github.com/hedzr/lb/lbapi"
	"hash/crc32"
	"sort"
	"sync"
)

// New make a new load-balancer instance with Ketama Hashing algorithm
func New(opts ...lbapi.Opt) lbapi.Balancer {
	return (&hashS{
		hasher:  crc32.ChecksumIEEE,
		replica: 32,
		keys:    make(map[uint32]lbapi.Peer),
		peers:   make(map[lbapi.Peer]bool),
	}).init(opts...)
}

// WithHashFunc allows a custom hash function to be specified.
// The default Hasher hash func is crc32.ChecksumIEEE.
func WithHashFunc(hashFunc Hasher) lbapi.Opt {
	return func(balancer lbapi.Balancer) {
		if l, ok := balancer.(*hashS); ok {
			l.hasher = hashFunc
		}
	}
}

// WithReplica allows a custom replica number to be specified.
// The default replica number is 32.
func WithReplica(replica int) lbapi.Opt {
	return func(balancer lbapi.Balancer) {
		if l, ok := balancer.(*hashS); ok {
			l.replica = replica
		}
	}
}

// Hasher is a hash function
type Hasher func(data []byte) uint32

// hashS is a impl with ketama consist hash algor
type hashS struct {
	hasher   Hasher
	replica  int
	hashRing []uint32
	keys     map[uint32]lbapi.Peer
	peers    map[lbapi.Peer]bool
	rw       sync.RWMutex
}

func (s *hashS) init(opts ...lbapi.Opt) *hashS {
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *hashS) Next(factor lbapi.Factor) (next lbapi.Peer, c lbapi.Constrainable) {
	var hash uint32
	if h, ok := factor.(lbapi.FactorHashable); ok {
		hash = h.HashCode()
	} else {
		hash = s.hasher([]byte(factor.Factor()))
	}

	next = s.miniNext(hash)
	if next != nil {
		if fc, ok := factor.(lbapi.FactorComparable); ok {
			next, c, _ = fc.ConstrainedBy(next)
		} else if nested, ok := next.(lbapi.BalancerLite); ok {
			next, c = nested.Next(factor)
		}
	}

	return
}

func (s *hashS) miniNext(hash uint32) (next lbapi.Peer) {
	s.rw.RLock()
	defer s.rw.RUnlock()

	ix := sort.Search(len(s.hashRing), func(i int) bool {
		return s.hashRing[i] >= hash
	})

	if ix == len(s.hashRing) {
		ix = 0
	}

	hashValue := s.hashRing[ix]
	if p, ok := s.keys[hashValue]; ok {
		if _, ok = s.peers[p]; ok {
			next = p
		}
	}

	return
}

func (s *hashS) Count() int {
	s.rw.RLock()
	defer s.rw.RUnlock()
	return len(s.peers)
}

func (s *hashS) Add(peers ...lbapi.Peer) {
	s.rw.Lock()
	defer s.rw.Unlock()

	for _, p := range peers {
		s.peers[p] = true
		for i := 0; i < s.replica; i++ {
			hash := s.hasher(s.peerToBinaryID(p, i))
			s.hashRing = append(s.hashRing, hash)
			s.keys[hash] = p
		}
	}

	sort.Slice(s.hashRing, func(i, j int) bool {
		return s.hashRing[i] < s.hashRing[j]
	})
}

func (s *hashS) peerToBinaryID(p lbapi.Peer, replica int) []byte {
	str := fmt.Sprintf("%v-%05d", p, replica)
	return []byte(str)
}

func (s *hashS) Remove(peer lbapi.Peer) {
	s.rw.Lock()
	defer s.rw.Unlock()

	if _, ok := s.peers[peer]; ok {
		delete(s.peers, peer)
	}

	var keys []uint32
	var km = make(map[uint32]bool)
	for i, p := range s.keys {
		if p == peer {
			keys = append(keys, i)
			km[i] = true
		}
	}

	for _, key := range keys {
		delete(s.keys, key)
	}

	var vn []uint32
	for _, x := range s.hashRing {
		if _, ok := km[x]; !ok {
			vn = append(vn, x)
		}
	}
	s.hashRing = vn
}

func (s *hashS) Clear() {
	s.rw.Lock()
	defer s.rw.Unlock()
	s.hashRing = nil
	s.keys = make(map[uint32]lbapi.Peer)
	s.peers = make(map[lbapi.Peer]bool)
}
