// Package server implements a very simple HTTP server that reports the last outcomes.
package server

import (
	"fmt"
	"net/http"
	"smon/logger"
)

// Start starts an HTTP server given a TCP port in a goroutine.
func Start(port int) {
	logger.Std.Infof("reporting server on port %v", port)
	if port > 0 {
		http.HandleFunc("/", handler)
		go http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<H1>smon results</H1><HR/><PRE>")
	for _, l := range logger.Std.BufferedLines() {
		fmt.Fprint(w, *l)
	}
	fmt.Fprint(w, "</PRE><HR>")
}
