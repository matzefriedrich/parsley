package types

import "errors"

const (
	ErrorRequiresFunctionValue = "the given value is not function"
)

var (
	ErrRequiresFunctionValue = errors.New(ErrorRequiresFunctionValue)
)

type registryError struct {
	FuncError
	serviceTypeName string
}

var _ FunqErrorWithServiceTypeName = &registryError{}

func (r *registryError) ServiceTypeName(name string) {
	r.serviceTypeName = name
}

var _ FunqErrorWithServiceTypeName = &registryError{}

func NewRegistryError(msg string, initializers ...FuncErrorFunc) error {
	err := &registryError{
		FuncError: FuncError{
			msg: msg,
		}}
	for _, f := range initializers {
		f(err)
	}
	return err
}
