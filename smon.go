// Package main runs the dummy checker.
package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"smon/checker"
	"smon/checker/dummy"
	"smon/checker/localip"
	"smon/checker/lookuphost"
	"smon/checker/traceroute"
	"smon/checker/wget"
	"smon/logger"
	"smon/monitor"
	"smon/reporter/logreport"
	"smon/reportstreamer"
	"smon/server"
	"time"
)

const usage = `

Usage: smon [--FLAGS] CONFIGFILE
To see useful flags, try smon --help.
The CONFIGFILE configures what checkers to run and with which intervals.
See smon.json.sample for an example.

`

var (
	port = flag.Int("p", 0, "port to start a reporting HTTP server, <1 for off")
)

type config struct {
	// JSON-held fields
	Checker    string `json:"checker"`
	Arg        string `json:"arg"`
	Interval   string `json:"interval"`
	MaxRuntime string `json:"maxruntime"`

	// Interval and maxruntime as a duration
	interval   time.Duration
	maxruntime time.Duration
}

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		logger.Std.Fatal(usage)
	}
	if *port > 0 {
		server.Start(*port)
	}

	logger.Std.Msg("starting")
	js, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		logger.Std.Fatalf("cannot read config file %v: %v", flag.Arg(0), err)
	}
	var checkers []config
	if err := json.Unmarshal(js, &checkers); err != nil {
		logger.Std.Fatalf("invalid JSON payload in config file %v: %v", flag.Arg(1), err)
	}

	mon := monitor.New(reportstreamer.New(logreport.New()))

	for _, c := range checkers {
		if c.Checker == "" {
			logger.Std.Fatal("configuration block lacks a checker ID")
		}
		if c.Interval == "" {
			logger.Std.Fatalf("checker %q lacks an interval configuration", c.Checker)
		}
		c.interval, err = time.ParseDuration(c.Interval)
		if err != nil {
			logger.Std.Fatalf("invalid duration %q in configuration: %v", c.Interval, err)
		}
		if c.MaxRuntime != "" {
			c.maxruntime, err = time.ParseDuration(c.MaxRuntime)
			if err != nil {
				logger.Std.Fatalf("invalid maxruntime %q in configuration: %v", c.MaxRuntime, err)
			}
		}

		var checker checker.Checker
		switch c.Checker {
		case "Dummy":
			checker = &dummy.Dummy{}
		case "LocalIP":
			checker = &localip.LocalIP{}
		case "LookupHost":
			checker = &lookuphost.LookupHost{}
		case "TraceRoute":
			checker = &traceroute.TraceRoute{}
		case "WGet":
			checker = &wget.WGet{}
		default:
			logger.Std.Fatalf("unknown checker %q", c.Checker)
		}
		if err := checker.Arg(c.Arg); err != nil {
			logger.Std.Fatalf("incorrect checker argument: %v", err)
		}
		if err := mon.Schedule(&monitor.Request{
			Checker:    checker,
			Interval:   c.interval,
			MaxRunTime: c.maxruntime,
		}); err != nil {
			logger.Std.Fatalf("failed to schedule checker %q: %v", c.Checker, err)
		}
	}
	mon.Run()
}
