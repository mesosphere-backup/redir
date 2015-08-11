package main

import (
	"fmt"
	"net/http"

	"github.com/miekg/dns"
)

// redirectHandler returns an http.Handler which resolves all requests
// of the form /path/:name and responds with an HTTP redirect
// pointing to the record chosen by the given strategy.
func redirectHandler(c client, path, resolver string, st strategy) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := dns.Fqdn(r.URL.Path[len(path):])

		var msg dns.Msg
		msg.SetQuestion(name, dns.TypeSRV)

		res, _, err := c.Exchange(&msg, resolver)
		if err != nil {
			code := http.StatusInternalServerError
			http.Error(w, err.Error(), code)
			return
		}

		srv := st(res.Answer)
		if srv == nil {
			code := http.StatusBadRequest
			http.Error(w, http.StatusText(code), code)
			return
		}

		location := fmt.Sprintf("http://%s:%d", srv.Target, srv.Port)
		http.Redirect(w, r, location, http.StatusSeeOther)
	})
}
