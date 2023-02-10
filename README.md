# go-lb

![Go](https://github.com/hedzr/lb/workflows/Go/badge.svg)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/hedzr/lb.svg?label=release)](https://github.com/hedzr/lb/releases)
[![](https://img.shields.io/badge/go-dev-green)](https://pkg.go.dev/github.com/hedzr/lb)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/hedzr/lb) <!-- [![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fhedzr%2Flb.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fhedzr%2Flb?ref=badge_shield) 
--> [![Go Report Card](https://goreportcard.com/badge/github.com/hedzr/lb)](https://goreportcard.com/report/github.com/hedzr/lb)
[![Coverage Status](https://coveralls.io/repos/github/hedzr/lb/badge.svg?branch=master&.9)](https://coveralls.io/github/hedzr/lb?branch=master)

`go-lb` provides a generic load balancers library.

## Features

The stocked algorithm are:

- random
- round-robin
- weighted round-robin
- consistent hash
- weighted random
- weighted versioning

Use `Register(...)/Unregister(...)` to add the balancer with your algorithm and use it with our `New(algorithm, opts...)`.

## History

- v0.5.0
  - a tiny logger interface has been embedded. So
    - we removed the dep to `hedzr/log` and free you from it
    - you may still enable internal logging sentences in `hedzr/lb` by setting up a custom [Logger](https://github.com/hedzr/lb/blob/master/pkg/logger/public.go#L4)

      ```go
      import "github.com/hedzr/lb/pkg/logger"
      logger.SetLogger(yoursLogger)
      ```

  - all codes reviewed

- v0.3.3
  - needs go modules 1.17 and higher
  - upgraded deps
  - remove unecessary deps
  - remove the nest go.mod since it cannot work any more
  - fix gin vuln report

- v0.3.1, v0.3.0 : work for below go1.17

## Usages

### Simple

```go
package main

import (
	"fmt"
	lb "github.com/hedzr/lb"
	"github.com/hedzr/lb/lbapi"
)

func main() {
	b := lb.New(lb.RoundRobin)

	b.Add(exP("172.16.0.7:3500"), exP("172.16.0.8:3500"), exP("172.16.0.9:3500"))
	sum := make(map[lbapi.Peer]int)
	for i := 0; i < 300; i++ {
		p, _ := b.Next(lbapi.DummyFactor)
		sum[p]++
	}

	for k, v := range sum {
		fmt.Printf("%v: %v\n", k, v)
	}
}

type exP string

func (s exP) String() string { return string(s) }
```

## About `Weighted versioning`

With the `Weighted versioning` algorithm, a set of version constraints and weights can be put as the basic rule. Such as:

1. "<= 1.1.x", weight: 2,
2. "^1.2.x", weight: 4,
3. "^2.x", weight: 11,
4. ^3.x", weight: 3,

the peers will be picked out by its version. And, all versioning peers will be picked with weights specified in constraints totally.

The partial codes:

```go
var testConstraints = []lbapi.Constrainable{
	version.NewConstrainablePeer("<= 1.1.x", 2),
	version.NewConstrainablePeer("^1.2.x", 4),
	version.NewConstrainablePeer("^2.x", 11),
	version.NewConstrainablePeer("^3.x", 3),
}

func TestVersionWRR2(t *testing.T) {
	lb := version.New(version.WithConstrainedPeers(testConstraints...))

	sum := make(map[lbapi.Peer]int)
	hits := make(map[lbapi.Peer]map[lbapi.Constrainable]bool)

	factor := initFactors()
	
	for i := 0; i < 500; i++ {
		peer, c := lb.Next(factor)

		sum[peer]++
		if ps, ok := hits[peer]; ok {
			if _, ok := ps[c]; !ok {
				ps[c] = true
			}
		} else {
			hits[peer] = make(map[lbapi.Constrainable]bool)
			hits[peer][c] = true
		}
	}

	// results
	total := 0
	for _, v := range sum {
		total += v
	}
	for peer, v := range sum {
		var keys []string
		var w int
		for c := range hits[peer] {
			if kk, ok := c.(fmt.Stringer); ok {
				keys = append(keys, kk.String())
			}
			if ww, ok := c.(lbapi.Weighted); ok {
				w = ww.Weight()
			}
		}
		// ex := findC(peer)
		t.Logf("%v: %v/%0.2f%%/w:%v. [%v => weight: %v]",
			peer, v, (float32(v)/float32(total))*100.0, w,
			strings.Join(keys, ","), w)
	}
}
```

See the full codes at [version/new_test.go](https://github.com/hedzr/lb/blob/master/version/new_test.go), 

For the full document of version constraints: [Masterminds/semver](https://github.com/Masterminds/semver) .


## API GW demo

Please check out the source codes:

- [gw-small](https://github.com/hedzr/lb/blob/master/_examples/gw-small/main.go)
- [svc-small](https://github.com/hedzr/lb/blob/master/_examples/svc-small/main.go)

And the command line to test its:

```bash
go run ./_examples/svc-small -port 8111 &
go run ./_examples/svc-small -port 8112 &
go run ./_examples/svc-small -port 8113 &

go run ./_examples/gw-lb-aware/main.go 8111 8112 8113 &

for ((i=0;i<5;i++)); do curl http://localhost:8103/ ; done
```



## License

MIT

