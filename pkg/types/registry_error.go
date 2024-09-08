package types

import "errors"

const (
	ErrorRequiresFunctionValue               = "the given value is not function"
	ErrorCannotRegisterModule                = "failed to register module"
	ErrorTypeAlreadyRegistered               = "type already registered"
	ErrorServiceAlreadyLinkedWithAnotherList = "service already linked with another list"
	ErrorFailedToRegisterType                = "failed to register type"
)

var (
	ErrRequiresFunctionValue = errors.New(ErrorRequiresFunctionValue)
	ErrCannotRegisterModule  = errors.New(ErrorCannotRegisterModule)
	ErrTypeAlreadyRegistered = errors.New(ErrorTypeAlreadyRegistered)
	ErrFailedToRegisterType  = errors.New(ErrorFailedToRegisterType)
)

type registryError struct {
	ParsleyError
	serviceTypeName string
}

var _ ParsleyErrorWithServiceTypeName = &registryError{}

func (r *registryError) ServiceTypeName(name string) {
	r.serviceTypeName = name
}

var _ ParsleyErrorWithServiceTypeName = &registryError{}

func NewRegistryError(msg string, initializers ...ParsleyErrorFunc) error {
	err := &registryError{
		ParsleyError: ParsleyError{
			Msg: msg,
		},
	}
	for _, f := range initializers {
		f(&err.ParsleyError)
		f(err)
	}
	return err
}
