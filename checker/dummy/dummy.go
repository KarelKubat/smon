// Package dummy implements a dummy service checker for testing.
package dummy

import (
	"fmt"
	"math/rand"
	"smon/checker"
	"smon/logger"
	"time"
)

// Dummy is the receiver for the dummy checker.
type Dummy struct {
	failureRate float32
	minRunTime  time.Duration
	maxRunTime  time.Duration
}

// Arg implements checker.Interface.
func (d *Dummy) Arg(s string) error {
	if s == "" {
		return fmt.Errorf(`Dummy checker requires a failure rate, a min runtime and a max runtime; e.g.: "0.24 2 4"`)
	}
	var frate float32
	var min, max int
	_, err := fmt.Sscanf(s, "%f %d %d", &frate, &min, &max)
	if err != nil {
		return fmt.Errorf("Dummy checker expects 3 numbers as its argument, not %q", s)
	}
	d.failureRate = frate
	d.minRunTime = time.Second * time.Duration(min)
	d.maxRunTime = time.Second * time.Duration(max)
	return nil
}

// Name implements checker.Interface.
func (d *Dummy) Name() string {
	return "Dummy"
}

// Run implements checker.Interface.
func (d *Dummy) Run() (checker.Outcome, error) {
	// Sleep for somewhere between the min and max runtime.
	diff := d.maxRunTime - d.minRunTime
	sleepDuration := time.Duration(rand.Int63n(int64(diff))) + d.minRunTime
	var outcome checker.Outcome

	if rand.Float32() < d.failureRate {
		outcome = checker.Failure
	} else {
		outcome = checker.Success
	}

	logger.Infof("Dummy: will return %v after %v", outcome, sleepDuration)
	time.Sleep(sleepDuration)
	return outcome, nil
}
