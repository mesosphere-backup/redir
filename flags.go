package main

import (
	"fmt"
	"net"
)

// address represents a "host/IP:port" pair and implements the flag.Value
// interface with fitting parsing and validation logic.
type address string

func (a address) String() string { return string(a) }

func (a *address) Set(val string) error {
	if _, port, err := net.SplitHostPort(val); err != nil {
		return fmt.Errorf("address is invalid: %s", err)
	} else if port == "" {
		return fmt.Errorf("address %q doesn't contain a port", val)
	}
	*a = address(val)
	return nil
}
