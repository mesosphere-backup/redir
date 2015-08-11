package main

import (
	"time"

	"github.com/miekg/dns"
)

// client defines an interface for a DNS client.
type client interface {
	Exchange(*dns.Msg, string) (*dns.Msg, time.Duration, error)
}

// newClient returns a dns.Client with the given timeout.
func newClient(timeout time.Duration) client {
	return &dns.Client{
		Net:            "udp",
		DialTimeout:    timeout,
		ReadTimeout:    timeout,
		WriteTimeout:   timeout,
		SingleInflight: true,
	}
}
