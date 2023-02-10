// Copyright © 2021 Hedzr Yeh.

package version

import (
	"testing"

	"github.com/hedzr/lb/rr"
)

func TestBackendFactor(t *testing.T) {
	bf := &backendFactor{
		version:    "3.1.9-rel",
		addr:       "abcd",
		versionObj: nil,
	}
	t.Logf("bf: %v, version = %v, factor = %v", bf, bf.Version(), bf.Factor())
}

func TestBackendsFactor(t *testing.T) {
	bf := &backendFactor{
		version:    "3.1.9-rel",
		addr:       "abcd",
		versionObj: nil,
	}
	t.Logf("bf: %v, version = %v, factor = %v", bf, bf.Version(), bf.Factor())

	bef := NewBackendsFactor(rr.New)
	bef.AddPeers(bf)
	t.Logf("bef: %v, factor = %v", bef, bef.Factor())
}
