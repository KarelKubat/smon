// Package logger implements logging for smon.
package logger

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	verbose = flag.Bool("v", false, "more logging")
)

// Error sends an error message to the logger.
func Error(s string) {
	log.SetOutput(os.Stderr)
	for _, part := range strings.Split(s, "\n") {
		log.Printf("ERROR: %v", part)
	}
}

// Errorf is like Error, but uses printf-like expansion.
func Errorf(f string, args ...interface{}) {
	Error(fmt.Sprintf(f, args...))
}

// Fatal sends an error message to the logger and halts the program.
func Fatal(s string) {
	log.SetOutput(os.Stderr)
	for _, part := range strings.Split(s, "\n") {
		log.Printf("FATAL: %v", part)
	}
	os.Exit(1)
}

// Fatalf is like Fatal but uses printf-like expansion.
func Fatalf(f string, args ...interface{}) {
	Fatal(fmt.Sprintf(f, args...))
}

// Info sends an informational message to the logger.
func Info(s string) {
	if !*verbose {
		return
	}
	log.SetOutput(os.Stdout)
	for _, part := range strings.Split(s, "\n") {
		log.Printf("INFO: %v", part)
	}
}

// Infof is like Info, but uses printf-like expansion.
func Infof(f string, args ...interface{}) {
	Info(fmt.Sprintf(f, args...))
}
