// Package reportstreamer handles incoming reports.
package reportstreamer

import (
	"smon/logger"
	"smon/monitor/result"
	"smon/reporter"
)

// ReportStreamer is the receiver for a streaming report generator.
type ReportStreamer struct {
	reporter reporter.Reporter
}

// New instatiates a new report streamer.
func New(r reporter.Reporter) *ReportStreamer {
	s := &ReportStreamer{
		reporter: r,
	}
	logger.Std.Info("reporstreamer: instantiated")
	return s
}

// Handle reads items from a stream of monitor results and reports them.
func (r *ReportStreamer) Handle(stream chan *result.Result) {
	for report := range stream {
		if err := r.reporter.Report(report); err != nil {
			logger.Std.Errorf("reportstreamer: %v", err)
		}
	}
}
