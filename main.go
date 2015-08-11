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
  This program starts a CORS enabled HTTP server on a given -addr whose
  requests to -path are converted into -timeout bound DNS SRV queries against
  the -resolver. The -param defines where the name to be resolved comes from. 
  Answers returned are load balanced with the given -strategy and used to respond
  to the original request with the defined -code and derived location.
`

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	cfg := config{addr: ":8080", resolver: ":53"}
	fs := flag.NewFlagSet("redir", flag.ContinueOnError)
	fs.StringVar(&cfg.path, "path", "/go/", "HTTP path to handle")
	fs.StringVar(&cfg.param, "param", "request-path", "DNS SRV name source [request-path, host-header]")
	fs.StringVar(&cfg.strategy, "strategy", "random", "SRV RR load balancing strategy [random, round-robin]")
	fs.StringVar(&cfg.origin, "origin", "*", "HTTP CORS Origin to accept")
	fs.IntVar(&cfg.code, "code", http.StatusTemporaryRedirect, "HTTP code to respond with")
	fs.DurationVar(&cfg.timeout, "timeout", time.Second, "DNS query timeout")
	fs.Var(&cfg.addr, "addr", "HTTP address to listen on")
	fs.Var(&cfg.resolver, "resolver", "DNS resolver addr to use")

	fs.SetOutput(os.Stderr)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "\nusage: redir [flags]\n")
		fs.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\ndescription:%s\n", description)
	}

	if err := fs.Parse(os.Args[1:]); err != nil {
		os.Exit(1)
	}

	strategy, err := cfg.Strategy()
	if err != nil {
		fatal(fs, err)
	}

	param, err := cfg.Param()
	if err != nil {
		fatal(fs, err)
	}

	cli := newClient(cfg.timeout)
	http.Handle(cfg.path, decorate(
		redirectHandler(cli, cfg.resolver.String(), cfg.code, strategy, param),
		methods("GET"),
		CORS(cfg.origin),
		logging(os.Stdout),
	))

	fmt.Fprintf(os.Stderr, "Listening on %q and redirecting %s requests\n", cfg.addr, cfg.path)
	if err := http.ListenAndServe(string(cfg.addr), nil); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// config holds the redir program configuration and exposes methods which
// validate it.
type config struct {
	path     string
	param    string
	strategy string
	origin   string
	code     int
	timeout  time.Duration
	addr     address
	resolver address
}

// Strategy returns the configured load balancing strategy function.
func (cfg config) Strategy() (strategy, error) {
	switch cfg.strategy {
	case "round-robin", "rr":
		return roundRobin(0), nil
	case "random", "rand":
		return random(time.Now().UnixNano()), nil
	default:
		return nil, fmt.Errorf("unsupported strategy %q", cfg.strategy)
	}
}

// Param returns the configured DNS SRV name param function.
func (cfg config) Param() (param, error) {
	switch cfg.param {
	case "request-path", "path":
		return path(cfg.path), nil
	case "host-header", "host":
		return header("Host"), nil
	default:
		return nil, fmt.Errorf("unsupported param %q", cfg.param)
	}
}

func fatal(fs *flag.FlagSet, err error) {
	fmt.Fprintln(os.Stderr, err)
	fs.Usage()
	os.Exit(1)
}
