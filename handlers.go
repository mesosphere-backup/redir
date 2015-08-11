package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/miekg/dns"
)

// redirectHandler returns an http.Handler which resolves all requests
// of the form /path/:name and responds with an HTTP with the given code,
// pointing to the record chosen by the given strategy.
func redirectHandler(c client, resolver string, code int, st strategy, name param) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var msg dns.Msg
		msg.SetQuestion(dns.Fqdn(name(r)), dns.TypeSRV)

		if res, _, err := c.Exchange(&msg, resolver); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else if srv := st(res.Answer); srv == nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		} else {
			location := fmt.Sprintf("http://%s:%d", srv.Target, srv.Port)
			http.Redirect(w, r, location, http.StatusSeeOther)
		}
	})
}

// param is a function type that extracts a string parameter from an *http.Request
type param func(*http.Request) string

// path returns a param function that strips the given prefix from an http.Request
// URL Path and returns the rest.
func path(strip string) param {
	return func(r *http.Request) string {
		return r.URL.Path[len(strip):]
	}
}

// header returns a param function that returns the given header from an http.Request
func header(name string) param {
	return func(r *http.Request) string {
		if strings.ToLower(name) == "host" {
			return r.Host
		}
		return r.Header.Get(name)
	}
}
