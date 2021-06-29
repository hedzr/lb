// Copyright © 2021 Hedzr Yeh.

package version

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/hedzr/lb/lbapi"
)

// NewConstrainablePeer creates a peer and it will be added to wrr balancer.
// the peer holds version constraints expression and weight.
//
// a version constraints expression has these forms:
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
//     https://github.com/Masterminds/semver
//
// See also NewBackendsFactor and New
func NewConstrainablePeer(constraints string, weight int) (peer lbapi.Constrainable) {
	vc, err := semver.NewConstraint(constraints)
	if err == nil {
		peer = NewConstrainablePeerFromObj(vc, weight)
	}
	return
}

func NewConstrainablePeerFromObj(constraints *semver.Constraints, weight int) (peer lbapi.Constrainable) {
	peer = &constrainablePeer{
		weight:         weight,
		constraints:    constraints.String(),
		constraintsObj: constraints,
	}
	return
}

type constrainablePeer struct {
	constraints    string
	weight         int
	constraintsObj *semver.Constraints
}

func (s *constrainablePeer) String() string                   { return s.constraints }
func (s *constrainablePeer) Weight() int                      { return s.weight }
func (s *constrainablePeer) Constraints() *semver.Constraints { return s.constraintsObj }
func (s *constrainablePeer) CanConstrain(o interface{}) (yes bool) {
	_, yes = o.(*semver.Version)
	return
}
func (s *constrainablePeer) Check(factor interface{}) (satisfied bool) {
	if s.constraintsObj == nil {
		var err error
		s.constraintsObj, err = semver.NewConstraint(s.constraints)
		if err != nil {
			fmt.Printf("illegal constraints: %q. %v\n", s.constraints, err)
		}
	}

	if v, ok := factor.(*semver.Version); ok {
		satisfied = s.constraintsObj.Check(v)
	} else if v, ok := factor.(interface{ Version() *semver.Version }); ok {
		satisfied = s.constraintsObj.Check(v.Version())
	}
	return
}
