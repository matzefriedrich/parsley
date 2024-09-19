package types

import "errors"

// ParsleyError represents an error with an associated message and optional underlying cause.
type ParsleyError struct {
	cause error
	Msg   string
}

// Error returns the message associated with the ParsleyError.
func (f ParsleyError) Error() string {
	return f.Msg
}

// Unwrap returns the underlying cause of the ParsleyError, allowing for error unwrapping functionality.
func (f ParsleyError) Unwrap() error {
	return f.cause
}

// Is compares the current ParsleyError's message with another error's message to determine if they are the same.
func (f ParsleyError) Is(err error) bool {
	return f.Error() == err.Error()
}

// ParsleyErrorFunc is a function type that modifies a given error.
type ParsleyErrorFunc func(e error)

// WithCause wraps a given error within a ParsleyError.
func WithCause(err error) ParsleyErrorFunc {
	return func(e error) {
		var funqErr *ParsleyError
		ok := errors.As(e, &funqErr)
		if ok {
			funqErr.cause = err
		}
	}
}

// ParsleyErrorWithServiceTypeName defines an interface for setting the service type name on errors.
type ParsleyErrorWithServiceTypeName interface {
	ServiceTypeName(name string)
}

// ForServiceType creates a ParsleyErrorFunc that sets the service type name on errors that implement the ParsleyErrorWithServiceTypeName interface.
func ForServiceType(serviceType string) ParsleyErrorFunc {
	return func(e error) {
		withServiceTypeErr, ok := e.(ParsleyErrorWithServiceTypeName)
		if ok {
			withServiceTypeErr.ServiceTypeName(serviceType)
		}
	}
}

// ParsleyAggregateError represents an aggregate of multiple errors.
type ParsleyAggregateError struct {
	errors []error
	Msg    string
}

// Error returns the message associated with the ParsleyAggregateError.
func (f ParsleyAggregateError) Error() string {
	return f.Msg
}

// Is checks if the given error is equivalent to any error within the ParsleyAggregateError.
func (f ParsleyAggregateError) Is(err error) bool {
	if f.Error() == err.Error() {
		return true
	}
	for _, cause := range f.errors {
		if errors.Is(err, cause) {
			return true
		}
	}
	return false
}

// WithAggregatedCause returns a ParsleyErrorFunc that sets an aggregated error cause with the provided errors.
func WithAggregatedCause(errs ...error) ParsleyErrorFunc {
	return func(e error) {
		var funqErr *ParsleyError
		ok := errors.As(e, &funqErr)
		if ok {
			funqErr.cause = &ParsleyAggregateError{
				errors: errs,
				Msg:    "one or more errors occurred",
			}
		}
	}
}
