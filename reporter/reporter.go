// Package reporter defines how a service monitor reporter must behave.
package reporter

import (
	"smon/monitor/result"
)

// Reporter is an interface that all reporters must adhere to.
type Reporter interface {
	Report(r *result.Result) error
}
