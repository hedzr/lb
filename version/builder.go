// Copyright © 2021 Hedzr Yeh.

package version

import (
	"github.com/Masterminds/semver/v3"
	"github.com/hedzr/lb/lbapi"
	"github.com/hedzr/lb/wrr"
	"github.com/hedzr/log"
)

// NewBuilder wraps an object into VerPeer struct.
//
// a verConstraints is an semver expression, such as:
//
//     ">= 1.2, < 3.0.0 || >= 4.2.3"
//     "<2.x"
//     "<3.1.x"
//     "~1.2.3"   // is equivalent to >= 1.2.3, < 1.3.0
//     "^1.2.3"   // is equivalent to >= 1.2.3, < 2.0.0
//     "1.2 - 1.4.5"
//     ...
//
// The basic comparisons are:
//
//     =: equal (aliased to no operator)
//     !=: not equal
//     >: greater than
//     <: less than
//     >=: greater than or equal to
//     <=: less than or equal to
//
// The advanced constraints expression syntax is described at:
//
//     https://github.com/Masterminds/semver/
//
func NewBuilder() Builder {
	return &builder{}
}

// Builder interface
type Builder interface {
	AddPeer(verConstraints string, weight int) Builder
	AddPeers(cs ...lbapi.Constrainable) Builder
	Build() lbapi.Balancer
}

type builder struct {
	peers []lbapi.Peer
}

func (b *builder) AddPeer(verConstraints string, weight int) Builder {
	vp := &constrainablePeer{
		constraints:    verConstraints,
		constraintsObj: nil,
		weight:         weight,
	}

	var err error
	vp.constraintsObj, err = semver.NewConstraint(vp.constraints)
	if err != nil {
		log.Errorf("version constraints %q parsing failed: %v", vp.constraints, err)
	}

	b.peers = append(b.peers, vp)
	return b
}

func (b *builder) AddPeers(cs ...lbapi.Constrainable) Builder {
	for _, vp := range cs {
		b.peers = append(b.peers, vp)
	}
	return b
}

func (b *builder) Build() lbapi.Balancer {
	lb := wrr.New()
	for _, vp := range b.peers {
		lb.Add(vp)
	}
	return lb
}

//type wrrPeer struct {
//	vc             string
//	verConstraints *semver.Constraints
//	weight         int
//}
//
//func (s *wrrPeer) VC() *semver.Constraints    { return s.verConstraints }
//func (s *wrrPeer) VersionConstraints() string { return s.vc }
//func (s *wrrPeer) String() string             { return s.vc }
//func (s *wrrPeer) Weight() int                { return s.weight }

//// VerPeer is a lbapi.Peer entity.
//type VerPeer struct {
//	weight         int
//	vc             string
//	verConstraints *semver.Constraints
//	object         interface{}
//}
//
//func (s *VerPeer) Payload() interface{}       { return s.object }
//func (s *VerPeer) VC() *semver.Constraints    { return s.verConstraints }
//func (s *VerPeer) VersionConstraints() string { return s.vc }
//func (s *VerPeer) String() string             { return s.vc }
//func (s *VerPeer) Weight() int                { return s.weight }
