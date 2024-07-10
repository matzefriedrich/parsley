package types

import "errors"

type FuncError struct {
	cause error
	msg   string
}

func (f FuncError) Error() string {
	return f.msg
}

func (f FuncError) Unwrap() error {
	return f.cause
}

func (f FuncError) Is(err error) bool {
	return f.Error() == err.Error()
}

type FuncErrorFunc func(e error)

func WithCause(err error) FuncErrorFunc {
	return func(e error) {
		var funqErr *FuncError
		ok := errors.As(e, &funqErr)
		if ok {
			funqErr.cause = err
		}
	}
}

type FunqErrorWithServiceTypeName interface {
	ServiceTypeName(name string)
}

func ForServiceType(serviceType string) FuncErrorFunc {
	return func(e error) {
		withServiceTypeErr, ok := e.(FunqErrorWithServiceTypeName)
		if ok {
			withServiceTypeErr.ServiceTypeName(serviceType)
		}
	}
}
