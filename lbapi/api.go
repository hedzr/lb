// Copyright © 2021 Hedzr Yeh.

package lbapi

import "reflect"

// Peer is a backend object, such as a host+port, a
// http/https url, or a constraint expression, and so on.
type Peer interface {
	// String will return the main identity string of a peer.
	//
	// You would have to implement only this one func to plug
	// your object into the architecture created by
	// lbapi.Balancer
	String() string
}

// Weighted object
type Weighted interface {
	Weight() int
}

// WeightedPeer is a Peer with Weight property.
type WeightedPeer interface {
	Peer
	Weighted
}

// Factor is a factor parameter for BalancerLite.Next.
//
// If you won't known what should be passed into
// BalancerLite.Next, send DummyFactor as it.
//
// But when you are constructing a complex lb system
// or a multiple-level one such as a weighted
// versioning load balancer. Take a look at our
// version.New and version.VersioningBackendFactor.
type Factor interface {
	Factor() string
}

// BalancerLite represents a generic load balancer.
type BalancerLite interface {
	// Next will return the next backend.
	//
	// For some ones like consistent-hashed balancer, Next needs
	// a factor as its param. It might be the requesting URL
	// typically.
	// For else, you can pass a lbapi.DummyFactor or just an
	// empty string as factor.
	Next(factor Factor) (next Peer, c Constrainable)
}

// Balancer represents a generic load balancer.
// For the real world, Balancer is a useful interface rather
// than BalancerLite.
type Balancer interface {
	BalancerLite
	Count() int
	Add(peers ...Peer)
	Remove(peer Peer)
	Clear()
}

// Opt is a type prototype for New Balancer
type Opt func(balancer Balancer)

// FactorComparable is a composite interface which assembly Factor and constraint comparing.
type FactorComparable interface {
	Factor
	ConstrainedBy(constraints interface{}) (peer Peer, c Constrainable, satisfied bool)
}

// FactorHashable is a composite interface which assembly Factor and hashing computer.
type FactorHashable interface {
	Factor
	HashCode() uint32
}

// FactorString is a string type, it implements Factor interface.
type FactorString string

// Factor function impl Factor interface
func (s FactorString) Factor() string { return string(s) }

// DummyFactor will be used in someone like you does not known
// what on earth should be passed into BalancerLite.Next(factor).
const DummyFactor FactorString = ""

// Constrainable is an object who can be applied onto a factor ( BalancerLite.Next(factor) )
type Constrainable interface {
	CanConstrain(o interface{}) (yes bool)
	Check(o interface{}) (satisfied bool)
	Peer
}

// WeightedConstrainable is an object who can be applied onto a factor ( BalancerLite.Next(factor) )
type WeightedConstrainable interface {
	Constrainable
	Weighted
}

// DeepEqualAware could be concreted by a Peer so you could
// customize how to compare two peers, avoid reflect.DeepEqual
// bypass.
type DeepEqualAware interface {
	DeepEqual(b Peer) bool
}

// DeepEqual will be used in Balancer, and a Peer can bypass
// reflect.DeepEqual by implementing DeepEqualAware interface.
func DeepEqual(a, b Peer) (yes bool) {
	if a == b {
		return true
	}

	if e, ok := a.(DeepEqualAware); ok {
		return e.DeepEqual(b)
	}

	return reflect.DeepEqual(a, b)
}

//// Sum sum
//func Sum(a []int, fn func(it int)) {
//	for _, it := range a {
//		fn(it)
//	}
//}
//
//// SumMapKeys sum
//func SumMapKeys(a map[int]string, fn func(it int)) {
//	for it := range a {
//		fn(it)
//	}
//}
//
//// SumMapValues sum
//func SumMapValues(a map[string]int, fn func(it int)) {
//	for _, it := range a {
//		fn(it)
//	}
//}
