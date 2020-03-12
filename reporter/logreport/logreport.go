// Package logreport implements a Reporter that writes to the standard Logger.
package logreport

import (
	"smon/logger"
	"smon/monitor/result"
)

// LogReport is the receiver.
type LogReport struct{}

// New returns a Reporter to write to the stated io.Writer.
func New() *LogReport {
	return &LogReport{}
}

// Report implements the interface Reporter.
func (l *LogReport) Report(r *result.Result) error {
	if r.CheckerErr != nil || r.SchedulerErr != nil {
		logger.Std.Error(r.String())
	} else {
		logger.Std.Msg(r.String())
	}
	return nil
}
