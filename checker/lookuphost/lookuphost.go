// Package lookuphost implements a checker to resolve DNS names.
package lookuphost

import (
	"fmt"
	"net"
	"smon/checker"
	"smon/logger"
)

// LookupHost is the receiver.
type LookupHost struct {
	hostname string
}

// Arg implements checker.Interface.
func (l *LookupHost) Arg(s string) error {
	if s == "" {
		return fmt.Errorf("LookupHost requires an argument (host to lookup)")
	}
	l.hostname = s
	return nil
}

// Name implements checker.Interface.
func (l *LookupHost) Name() string {
	return fmt.Sprintf("LookupHost %s", l.hostname)
}

// Run implements checker.Interface.
func (l *LookupHost) Run() (checker.Outcome, error) {
	addrs, err := net.LookupHost(l.hostname)
	if err == nil {
		logger.Infof("LookupHost: %s -> %v", l.hostname, addrs)
		return checker.Success, nil
	}
	return checker.Failure, err
}
