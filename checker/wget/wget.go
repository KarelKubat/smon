// Package wget implements a checker to retrieve a web page.
package wget

import (
	"fmt"
	"net/http"
	"smon/checker"
	"smon/logger"
)

// WGet is the receiver.
type WGet struct {
	URL string
}

// Arg implements checker.Interface.
func (w *WGet) Arg(s string) error {
	if s == "" {
		return fmt.Errorf("WGet requires an argument (URL to fetch)")
	}
	w.URL = s
	return nil
}

// Name implements checker.Interface.
func (w *WGet) Name() string {
	return fmt.Sprintf("WGet %s", w.URL)
}

// Run implements checker.Interface.
func (w *WGet) Run() (checker.Outcome, error) {
	resp, err := http.Get(w.URL)
	if err != nil {
		return checker.Failure, err
	}
	defer resp.Body.Close()
	logger.Infof("WGet: %s succeeded", w.URL)
	return checker.Success, nil
}
