// Package checker defines how service checkers must behave.
package checker

// Outcome is the result of a service checker.
type Outcome int

const (
	// Undefined should never be returned, it represents a non-set value.
	Undefined Outcome = iota
	// Failure represents service unavailability.
	Failure
	// Success is service availablility.
	Success
)

// String returns a string representation of an outcome.
func (o Outcome) String() string {
	return []string{
		"Undefined",
		"Failure",
		"Success",
	}[int(o)]
}

// Checker defines the methods that a service checker must provide.
type Checker interface {
	Arg(string) error      // parses any arg from the config
	Name() string          // checker name
	Run() (Outcome, error) // perform the check
}
