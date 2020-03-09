// Package localip implements a checker that verifies that a local IP address exists.
package localip

import (
	"fmt"
	"net"
	"smon/checker"
	"smon/logger"
)

// LocalIP is the receiver.
type LocalIP struct{}

// Arg implements checker.Interface.
func (l *LocalIP) Arg(s string) error {
	if s != "" {
		return fmt.Errorf("LocalIP does not support an argument")
	}
	return nil
}

// Name implements checker.Interface.
func (l *LocalIP) Name() string {
	return "LocalIP"
}

// Run implements checker.Interface.
func (l *LocalIP) Run() (checker.Outcome, error) {
	// This doesn't actually connect, just checks from where we would dial.
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return checker.Failure, err
	}
	defer conn.Close()

	addr := conn.LocalAddr().(*net.UDPAddr)
	logger.Infof("LocalIP: %v", addr)
	return checker.Success, nil
}
