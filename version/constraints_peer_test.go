// Copyright © 2021 Hedzr Yeh.

package version

import (
	"testing"

	"github.com/Masterminds/semver/v3"
)

func TestNewConstraintsPeer(t *testing.T) {
	peer := NewConstrainablePeer("<= 1.1.x", 2)
	t.Logf("%v, %v, %v", peer, peer.CanConstrain(nil), peer.(*constrainablePeer).Constraints())

	x := &xS{ver: "1.1.13"}
	peer.Check(x)
}

func TestNewConstraintsPeer2(t *testing.T) {
	peer := &constrainablePeer{
		constraints:    "<= 1.1.x",
		weight:         2,
		constraintsObj: nil,
	}
	t.Logf("%v, %v, %v", peer, peer.CanConstrain(nil), peer.Constraints())

	x := &xS{ver: "1.1.13"}
	peer.Check(x)
}

type xS struct {
	ver string
}

func (s *xS) Version() *semver.Version {
	vs, _ := semver.NewVersion(s.ver)
	return vs
}
