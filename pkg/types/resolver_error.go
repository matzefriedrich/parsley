package types

import "errors"

const (
	ErrorServiceTypeNotRegistered     = "service type is not registered"
	ErrorRequiredServiceNotRegistered = "required service type is not registered"
	ErrorCannotResolveService         = "cannot resolve service"
)

var (
	ErrServiceTypeNotRegistered = errors.New(ErrorServiceTypeNotRegistered)
)

type ResolverError struct {
	FuncError
	serviceTypeName string
}

var _ FunqErrorWithServiceTypeName = &ResolverError{}

func (r *ResolverError) ServiceTypeName(name string) {
	r.serviceTypeName = name
}

func NewResolverError(msg string, initializers ...FuncErrorFunc) error {
	err := &ResolverError{
		FuncError: FuncError{
			msg: msg,
		},
	}
	for _, f := range initializers {
		f(err)
	}
	return err
}
