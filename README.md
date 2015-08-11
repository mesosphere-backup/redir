# Redir

Redir is an HTTP redirector of DNS SRV records with configurable
load-balancing strategies.

## Installing 
```shell
$ go get github.com/mesosphere/redir
```

## Running
```shell
$ redir -h

usage: redir [flags]
  -addr=:8080: HTTP address to listen on
  -origin="*": HTTP CORS Origin to accept
  -path="/go/": HTTP prefix path route to redirect
  -resolver=127.0.0.1:53: DNS resolver addr to use
  -strategy="random": SRV RR load balancing strategy [random, round-robin]
  -timeout=1s: DNS query timeout

description:
  This program starts an HTTP server on a given -addr whose requests to /:path/:name
  are converted into -timeout bound DNS SRV queries against the -resolver.
  Answers returned are load balanced with the given -strategy and used to respond
  the original request with an HTTP redirect to the chosen SRV record.

```
