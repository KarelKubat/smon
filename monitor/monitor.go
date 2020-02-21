// Package monitor schedules and executes service checkers
package monitor

import (
	"flag"
	"fmt"
	"smon/checker"
	"smon/logger"
	"smon/monitor/result"
	"smon/reportstreamer"
	"time"
)

var (
	continuous = flag.Bool("c", false,
		"continuously report checkers' outcomes (default: only at state change)")
)

// Request represents how a checker should be run.
type Request struct {
	// Checker is the checker to run.
	Checker checker.Checker
	// Interval is the time between invocations.
	Interval time.Duration
	// MaxRunTime is the maximum time that a checker is allowed to run. Ut may not be shorter than
	// the Interval.
	MaxRunTime time.Duration
}

// Monitor is a receiver where checkers are run.
type Monitor struct {
	reportChannel  chan *result.Result
	reportStreamer *reportstreamer.ReportStreamer
}

// New instantiates a Monitor.
func New(streamer *reportstreamer.ReportStreamer) *Monitor {
	m := &Monitor{
		reportChannel:  make(chan *result.Result),
		reportStreamer: streamer,
	}
	logger.Info("monitor: instantiated")
	return m
}

// ReportChannel returns the channel where checkers' results are sent.
func (m *Monitor) ReportChannel() chan *result.Result {
	return m.reportChannel
}

// Schedule starts asynchronously executing a checker.
func (m *Monitor) Schedule(req *Request) error {
	if req.MaxRunTime == 0 {
		req.MaxRunTime = req.Interval
	}
	if req.MaxRunTime < req.Interval {
		return fmt.Errorf("%q: scheduling request has maxruntime %v, which is shorter than the interval %v",
			req.Checker.Name(), req.MaxRunTime, req.Interval)
	}

	logger.Infof("monitor: scheduling %q", req.Checker.Name())
	go m.start(req)
	return nil
}

type checkerResult struct {
	runtime time.Duration
	outcome checker.Outcome
	err     error
}

// start is a helper that's started asynchronously.
func (m *Monitor) start(req *Request) {
	lastOutcome := checker.Undefined
	for {
		// Run the checker, send its output to the result stream.
		logger.Infof("monitor: invoking %q", req.Checker.Name())

		resultChan := make(chan checkerResult)
		timeoutChan := time.NewTimer(req.MaxRunTime)
		defer timeoutChan.Stop()

		start := time.Now()
		go func(chan checkerResult) {
			outcome, err := req.Checker.Run()
			runtime := time.Now().Sub(start)
			defer func() {
				// resultChan may be closed when the checkerResult is sent, causing a panic.
				if r := recover(); r != nil {
					logger.Errorf("%v exceeded maxruntime", req.Checker.Name())
				}
			}()
			resultChan <- checkerResult{
				runtime: runtime,
				outcome: outcome,
				err:     err,
			}
		}(resultChan)

		select {
		case <-timeoutChan.C:
			close(resultChan)
			m.reportChannel <- &result.Result{
				Name:         req.Checker.Name(),
				Start:        start,
				SchedulerErr: fmt.Errorf("maxruntime %v exceeded", req.MaxRunTime),
			}
			lastOutcome = checker.Undefined
		case o := <-resultChan:
			timeoutChan.Stop()
			if *continuous || lastOutcome != o.outcome {
				m.reportChannel <- &result.Result{
					Name:       req.Checker.Name(),
					Start:      start,
					Duration:   o.runtime,
					Outcome:    o.outcome,
					CheckerErr: o.err,
				}
			}
			lastOutcome = o.outcome
			time.Sleep(req.Interval - o.runtime)
		}
	}
}

// Run listens for checker results. It doesn't return.
func (m *Monitor) Run() {
	m.reportStreamer.Handle(m.reportChannel)
}
