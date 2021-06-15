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
