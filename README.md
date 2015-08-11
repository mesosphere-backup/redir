# Redir

Redir is an HTTP redirector of DNS SRV records with configurable
load-balancing strategies.

## Installing 

Pre-compiled binaries are [available](https://github.com/mesosphere/redir/releases)
for most OS and architectures. 

Alternatively, if you have Go installed:
```shell
$ go get github.com/mesosphere/redir
```

## Running
```shell
$ redir -h

usage: redir [flags]
  -addr=:8080: HTTP address to listen on
  -code=307: HTTP code to respond with
  -origin="*": HTTP CORS Origin to accept
  -param="request-path": DNS SRV name source [request-path, host-header]
  -path="/go/": HTTP path to handle
  -resolver=:53: DNS resolver addr to use
  -strategy="random": SRV RR load balancing strategy [random, round-robin]
  -timeout=1s: DNS query timeout

description:
  This program starts a CORS enabled HTTP server on a given -addr whose
  requests to -path are converted into -timeout bound DNS SRV queries against
  the -resolver. The -param defines where the name to be resolved comes from.
  Answers returned are load balanced with the given -strategy and used to respond
  to the original request with the defined -code and derived location.
```
