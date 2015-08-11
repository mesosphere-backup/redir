package main

import (
	"math/rand"
	"sync/atomic"

	"github.com/miekg/dns"
)

// strategy represents a load balancing strategy
type strategy func([]dns.RR) *dns.SRV

// random returns a random load balancing strategy seeded with the given seed.
func random(seed int64) strategy {
	rnd := rand.New(rand.NewSource(seed))
	return func(rrs []dns.RR) *dns.SRV {
		if srvs := collect(rrs); len(srvs) > 0 {
			return srvs[rnd.Intn(len(srvs))]
		}
		return nil
	}
}

// roundRobin returns a round-robin load balancing strategy starting with the
// given robin, going through the priority sorted SRV records.
func roundRobin(robin uint64) strategy {
	return func(rrs []dns.RR) *dns.SRV {
		if srvs := collect(rrs); len(srvs) > 0 {
			return srvs[atomic.AddUint64(&robin, 1)%uint64(len(srvs))]
		}
		return nil
	}
}

// collect collects SRV records out of a dns.RR slice.
func collect(rrs []dns.RR) []*dns.SRV {
	srvs := make([]*dns.SRV, 0, len(rrs))
	for _, r := range rrs {
		if srv, ok := r.(*dns.SRV); ok {
			srvs = append(srvs, srv)
		}
	}
	return srvs
}
