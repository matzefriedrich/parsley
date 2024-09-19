package types

import "errors"

const (
	ErrorInstanceAlreadySet = "instance already set"
)

var (
	// ErrInstanceAlreadySet is returned when there is an attempt to set an instance that is already set.
	ErrInstanceAlreadySet = errors.New(ErrorInstanceAlreadySet)
)

// DependencyError represents an error that occurs due to a missing or failed dependency.
// This error type encapsulates a ParsleyError.
type DependencyError struct {
	ParsleyError
}

// NewDependencyError creates a new DependencyError with the provided message.
func NewDependencyError(msg string) error {
	err := DependencyError{ParsleyError{Msg: msg}}
	return err
}
