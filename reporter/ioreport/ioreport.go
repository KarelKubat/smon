// Package ioreport implements a Reporter that writes to any io.Writer.
package ioreport

import (
	"fmt"
	"io"
	"smon/logger"
	"smon/monitor/result"
)

// IOReport is the receiver to send reports to an iowriter.
type IOReport struct {
	writer io.Writer
}

// New returns a Reporter to write to the stated io.Writer.
func New(w io.Writer) *IOReport {
	i := &IOReport{
		writer: w,
	}
	logger.Info("ioreport: instantiated")
	return i
}

// Report implements the interface Reporter.
func (iow *IOReport) Report(r *result.Result) error {
	fmt.Fprintf(iow.writer, "%v\n", r)
	return nil
}
