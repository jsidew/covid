package errors

import "fmt"

// Prefix for the errors created with F and W.
// The colon ":" between the prefix and the error message will be added automatically.
var Prefix string

type errwrap struct {
	pfx string
	err error
}

func (e *errwrap) Error() string { return e.pfx + ": " + e.err.Error() }
func (e *errwrap) Unwrap() error { return e.err }

// F returns a formatted error with Prefix
func F(format string, a ...interface{}) error {
	return fmt.Errorf(Prefix+": "+format, a...)
}

// W wraps an error with Prefix
func W(err error) error {
	return &errwrap{Prefix, err}
}
