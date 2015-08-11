package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"
)

const description = `
  This program starts an HTTP server on a given -addr whose requests to /:path/:name
  are converted into -timeout bound DNS SRV queries against the -resolver.
  Answers returned are load balanced with the given -strategy and used to respond
  the original request with an HTTP redirect to the chosen SRV record.
`

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var (
		addr     = address(":8080")
		resolver = address("127.0.0.1:53")
	)

	fs := flag.NewFlagSet("redir", flag.ContinueOnError)
	path := fs.String("path", "/go/", "HTTP prefix path route to redirect")
	strategy := fs.String("strategy", "random", "SRV RR load balancing strategy [random, round-robin]")
	origin := fs.String("origin", "*", "HTTP CORS Origin to accept")
	timeout := fs.Duration("timeout", time.Second, "DNS query timeout")
	fs.Var(&addr, "addr", "HTTP address to listen on")
	fs.Var(&resolver, "resolver", "DNS resolver addr to use")

	fs.SetOutput(os.Stderr)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "\nusage: redir [flags]\n")
		fs.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\ndescription:%s\n", description)
	}

	if err := fs.Parse(os.Args[1:]); err != nil {
		os.Exit(1)
	}

	st, ok := strategies[*strategy]
	if !ok {
		fmt.Fprintf(os.Stderr, "strategy %q is invalid", strategy)
		fs.Usage()
		os.Exit(1)
	}

	http.Handle(*path, decorate(
		redirectHandler(newClient(*timeout), *path, resolver.String(), st),
		methods("GET"),
		CORS(*origin),
		logging(os.Stdout),
	))

	fmt.Fprintf(os.Stderr, "Listening on %q and redirecting %s requests\n", addr, *path)
	if err := http.ListenAndServe(string(addr), nil); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
