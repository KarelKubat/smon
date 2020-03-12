// Package logger implements logging for smon.
package logger

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"smon/logger/store"
	"strings"
)

const (
	infoPrefix  = "INFO  "
	errorPrefix = "ERROR "
	fatalPrefix = "FATAL "
)

var (
	verbose     = flag.Bool("v", false, "more logging")
	linesToKeep = flag.Int("l", 512,
		"log entries to keep in memory for reporting (useful with reporting server)")

	// Std is a default Logger to send its output to stdout/err.
	Std = New(os.Stdout, os.Stderr)
)

// Logger is the receiver.
type Logger struct {
	buffered    *store.Store
	infoWriter  io.Writer
	errorWriter io.Writer
}

// New instantiates a Logger given writers for debug and error messaging.
func New(i, e io.Writer) *Logger {
	return &Logger{
		buffered:    store.New(*linesToKeep),
		infoWriter:  i,
		errorWriter: e,
	}
}

// Error sends an error message to the logger.
func (l *Logger) Error(s string) {
	l.output(l.errorWriter, errorPrefix, s)
}

// Errorf is like Error, but uses printf-like expansion.
func (l *Logger) Errorf(f string, args ...interface{}) {
	l.Error(fmt.Sprintf(f, args...))
}

// Fatal sends an error message to the logger and halts the program.
func (l *Logger) Fatal(s string) {
	l.output(l.errorWriter, fatalPrefix, s)
	os.Exit(1)
}

// Fatalf is like Fatal but uses printf-like expansion.
func (l *Logger) Fatalf(f string, args ...interface{}) {
	l.Fatal(fmt.Sprintf(f, args...))
}

// Info sends an informational message to the logger when verbosity is turned on.
func (l *Logger) Info(s string) {
	if !*verbose {
		return
	}
	l.Msg(s)
}

// Infof is like Info, but uses printf-like expansion.
func (l *Logger) Infof(f string, args ...interface{}) {
	l.Info(fmt.Sprintf(f, args...))
}

// Msg sends an informational message, regardless of verbosity.
func (l *Logger) Msg(s string) {
	l.output(l.infoWriter, infoPrefix, s)
}

// Msgf is like Msg, but uses printf-like expansion.
func (l *Logger) Msgf(f string, args ...interface{}) {
	l.Msg(fmt.Sprintf(f, args...))
}

// BufferedLines returns the lines that were stored up until now.
func (l *Logger) BufferedLines() []*string {
	return l.buffered.Lines()
}

// output is a helper function.
func (l *Logger) output(w io.Writer, prefix, s string) {
	for _, part := range strings.Split(s, "\n") {
		if part == "" {
			continue
		}
		// Log output to a local writer. Duplicate to the intended output stream and store.
		var buf bytes.Buffer
		log.New(&buf, prefix, log.Ldate|log.Ltime).Print(part)
		outString := buf.String()
		fmt.Fprint(w, outString)
		l.buffered.Add(outString)
	}
}
