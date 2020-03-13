// Package traceroute implements a checker to trace a route to a host.
package traceroute

import (
	"bytes"
	"fmt"
	"os/exec"
	"smon/checker"
	"smon/logger"
	"strings"
)

const (
	traceroute = "traceroute"
)

// TraceRoute is the receiver.
type TraceRoute struct {
	hostname string
}

// Arg implements checker.Interface.
func (t *TraceRoute) Arg(s string) error {
	if s == "" {
		return fmt.Errorf("TraceRoute requires an argument (hostname)")
	}
	t.hostname = s
	return nil
}

// Name implements checker.Interface.
func (t *TraceRoute) Name() string {
	return fmt.Sprintf("TraceRoute %s", t.hostname)
}

// collect is a helper to make stdout/stderr outputs more readable.
func collect(stdout, stderr string) string {
	report := ""
	if stdout != "" {
		report = fmt.Sprintf("stdout: [%s]", strings.Join(strings.Split(strings.TrimSuffix(stdout, "\n"), "\n"), " "))
	}
	if stderr != "" {
		if report != "" {
			report += ", "
		}
		report += fmt.Sprintf("stderr: [%s]", strings.Join(strings.Split(strings.TrimSuffix(stderr, "\n"), "\n"), " "))
	}
	return report
}

// Run implements checker.Interface.
func (t *TraceRoute) Run() (checker.Outcome, error) {
	_, err := exec.LookPath(traceroute)
	if err != nil {
		return checker.Failure, fmt.Errorf("TraceRoute: can't find %q along the $PATH", traceroute)
	}
	cmd := exec.Command(traceroute, t.hostname)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return checker.Failure, fmt.Errorf("TraceRoute: failed with %v, %s", err, collect(stdout.String(), stderr.String()))
	}
	logger.Std.Infof("Traceroute: %s -> %s", t.hostname, stdout.String())
	return checker.Success, nil
}
