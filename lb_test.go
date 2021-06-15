package lb_test

import (
	lb2 "github.com/hedzr/lb"
	"github.com/hedzr/lb/lbapi"
	"testing"
)

type exP struct {
	addr   string
	weight int
}

func (s *exP) String() string { return s.addr }
func (s *exP) Weight() int    { return s.weight }

func TestNew(t *testing.T) {
	lb := lb2.New(lb2.WeightedRoundRobin,
		lb2.WithPeers(&exP{"172.16.0.7:3500", 5}, &exP{"172.16.0.8:3500", 3}, &exP{"172.16.0.9:3500", 2}))

	sum := make(map[lbapi.Peer]int)

	for i := 0; i < 300; i++ {
		peer, _ := lb.Next(lbapi.DummyFactor)
		sum[peer]++
	}

	// results
	total := 0
	for _, v := range sum {
		total += v
	}
	for k, v := range sum {
		t.Logf("%v: %v. weight = %v, %0.2f%%", k, v, k.(lbapi.WeightedPeer).Weight(), (float32(v)/float32(total))*100.0)
	}
}

func Test_AddRemove(t *testing.T) {
	lb := lb2.New(lb2.WeightedRoundRobin)
	lb.Add(&exP{"172.16.0.7:3500", 5}, &exP{"172.16.0.8:3500", 3}, &exP{"172.16.0.9:3500", 2})

	lb.Add(&exP{"172.16.0.8:3500", 3})
	if lb.Count() != 3 {
		t.Fatal("wrong Add: the dup peer should be ignore")
	}

	lb.Remove(&exP{"172.16.0.8:3500", 3})
	if lb.Count() != 2 {
		t.Fatalf("wrong Remove: not removed? count = %v", lb.Count())
	}
}

func Test_Register(t *testing.T) {
	lb2.Register("nil", func(opts ...lbapi.Opt) lbapi.Balancer {
		return nil
	})

	lb2.New("nil")

	lb2.Unregister("nil")
}
