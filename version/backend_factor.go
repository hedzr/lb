package version

import (
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/hedzr/lb/lbapi"
	"github.com/hedzr/log"
)

// VersioningBackendFactor interface
type VersioningBackendFactor interface {
	lbapi.Factor
	lbapi.Peer
	Version() *semver.Version
}

// NewBackendFactor make a instance with (backend address, version) pair.
// version can be 'v1.2.3' or '1.2.3', 'v' will be striped.
// address is like 'host:port' typically, but you can use any forms you like.
func NewBackendFactor(version string, addr string) VersioningBackendFactor {
	f := &backendFactor{
		version: check(version),
		addr:    addr,
	}

	var err error
	f.versionObj, err = semver.NewVersion(f.version)
	if err != nil {
		log.Errorf("illegal version %q: %v", f.version, err)
	}

	return f
}

type backendFactor struct {
	version    string
	addr       string
	versionObj *semver.Version
}

func check(v string) string {
	if v[0] == 'v' {
		return v[1:]
	}
	return v
}

func (f *backendFactor) Factor() string { return f.version }
func (f *backendFactor) String() string { return fmt.Sprintf("%v - %v", f.addr, f.version) }
func (f *backendFactor) Version() *semver.Version {
	if f.versionObj == nil {
		f.versionObj, _ = semver.NewVersion(f.version)
	}
	return f.versionObj
}
