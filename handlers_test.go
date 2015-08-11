package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/miekg/dns"
)

func TestRedirectHandler(t *testing.T) {
	c := clientFunc(func(m *dns.Msg, _ string) (*dns.Msg, time.Duration, error) {
		var r dns.Msg
		r.SetReply(m)
		r.Answer = map[string][]dns.RR{
			"_abc._tcp.domain.": []dns.RR{
				&dns.SRV{Target: "abc.domain.", Port: 5000},
				&dns.SRV{Target: "abc.domain.", Port: 5001},
				&dns.SRV{Target: "abc.domain.", Port: 5002},
				&dns.SRV{Target: "abc.domain.", Port: 5003},
			},
		}[m.Question[0].Name]
		return &r, 0, nil
	})

	for i, tt := range []struct {
		strategy
		name     string
		code     int
		location string
	}{
		{roundRobin(0), "_abc._tcp.domain", 303, "http://abc.domain.:5001"},
		{roundRobin(1), "_abc._tcp.domain", 303, "http://abc.domain.:5002"},
		{roundRobin(2), "_abc._tcp.domain.", 303, "http://abc.domain.:5003"},
		{roundRobin(3), "_abc._tcp.domain.", 303, "http://abc.domain.:5000"},

		{random(0), "_abc._tcp.domain", 303, "http://abc.domain.:5002"},
		{random(1), "_abc._tcp.domain", 303, "http://abc.domain.:5001"},
		{random(2), "_abc._tcp.domain.", 303, "http://abc.domain.:5002"},
		{random(3), "_abc._tcp.domain.", 303, "http://abc.domain.:5000"},

		{roundRobin(0), "_missing._tcp.domain.", 400, ""},
		{random(0), "_missing._tcp.domain.", 400, ""},
	} {
		req, err := http.NewRequest("GET", "/"+tt.name, nil)
		if err != nil {
			t.Fatal(err)
		}
		rw := httptest.NewRecorder()
		h := redirectHandler(c, "/", "", tt.strategy)

		h.ServeHTTP(rw, req)
		if got, want := rw.Code, tt.code; got != want {
			t.Errorf("test #%d: got code: %d, want: %d", i, got, want)
		} else if got, want := rw.HeaderMap.Get("Location"), tt.location; got != want {
			t.Errorf("test #%d: got location: %v, want: %v", i, got, want)
		}
	}
}

// clientFunc is a function type that implements the client interface
type clientFunc func(*dns.Msg, string) (*dns.Msg, time.Duration, error)

func (f clientFunc) Exchange(m *dns.Msg, addr string) (*dns.Msg, time.Duration, error) {
	return f(m, addr)
}
