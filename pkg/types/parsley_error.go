package types

import "errors"

type ParsleyError struct {
	cause error
	msg   string
}

func (f ParsleyError) Error() string {
	return f.msg
}

func (f ParsleyError) Unwrap() error {
	return f.cause
}

func (f ParsleyError) Is(err error) bool {
	return f.Error() == err.Error()
}

type ParsleyErrorFunc func(e error)

func WithCause(err error) ParsleyErrorFunc {
	return func(e error) {
		var funqErr *ParsleyError
		ok := errors.As(e, &funqErr)
		if ok {
			funqErr.cause = err
		}
	}
}

type ParsleyErrorWithServiceTypeName interface {
	ServiceTypeName(name string)
}

func ForServiceType(serviceType string) ParsleyErrorFunc {
	return func(e error) {
		withServiceTypeErr, ok := e.(ParsleyErrorWithServiceTypeName)
		if ok {
			withServiceTypeErr.ServiceTypeName(serviceType)
		}
	}
}

type ParsleyAggregateError struct {
	errors []error
	msg    string
}

func (f ParsleyAggregateError) Error() string {
	return f.msg
}

func (f ParsleyAggregateError) Is(err error) bool {
	if f.Error() == err.Error() {
		return true
	}
	for _, err := range f.errors {
		if errors.Is(err, err) {
			return true
		}
	}
	return false
}

func WithAggregatedCause(errs ...error) ParsleyErrorFunc {
	return func(e error) {
		var funqErr *ParsleyError
		ok := errors.As(e, &funqErr)
		if ok {
			funqErr.cause = &ParsleyAggregateError{
				errors: errs,
				msg:    "one or more errors occurred",
			}
		}
	}
}
