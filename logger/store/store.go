// Package store keeps a buffer of the last log messages for display in the reporting server.
package store

// Store is the receiver.
type Store struct {
	nlines int
	lines  []*string
}

// New initializes a Store.
func New(n int) *Store {
	return &Store{
		nlines: n,
		lines:  []*string{},
	}
}

// Add adds a string to the Store to the front of the storage.
func (s *Store) Add(str string) {
	if len(s.lines) < s.nlines {
		s.lines = append(s.lines, &str)
	} else {
		s.lines = s.lines[1:s.nlines]
		s.lines[0] = &str
	}
}

// Lines returns the storage.
func (s *Store) Lines() []*string {
	return s.lines
}
