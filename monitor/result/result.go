// Package result wraps a checker's result and a monitor state.
package result

import (
	"fmt"
	"smon/checker"
	"time"
)

// Result represents the outcome of one checker run: its duration, the outcome of the checker,
// the error state of the checker, and the error state of the scheduler.
type Result struct {
	Name         string
	Start        time.Time
	Duration     time.Duration
	Outcome      checker.Outcome
	CheckerErr   error
	SchedulerErr error
}

// String returns a string representation of a Result.
func (r *Result) String() string {
	s := r.Start.Format("2006/01/02 15:04:05") + fmt.Sprintf(" %q: ", r.Name)
	switch {
	case r.SchedulerErr != nil:
		return s + fmt.Sprintf("SCHEDULER: %v", r.SchedulerErr)
	case r.CheckerErr != nil:
		return s + fmt.Sprintf("FAILURE: checker: %v", r.CheckerErr)
	default:
		return s + fmt.Sprintf("%v (after %v)", r.Outcome, r.Duration)
	}
}
