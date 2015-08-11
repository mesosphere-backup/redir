package main

import (
	"fmt"
	"net/http"

	"github.com/miekg/dns"
)

// redirectHandler returns an http.Handler which resolves all requests
// of the form /path/:name and responds with an HTTP with the given code,
// pointing to the record chosen by the given strategy.
func redirectHandler(c client, path, resolver string, code int, st strategy) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := dns.Fqdn(r.URL.Path[len(path):])
		msg := dns.Msg{}
		msg.SetQuestion(name, dns.TypeSRV)

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
