package types

import (
	"errors"
	"fmt"
)

const (
	ErrorRequiresFunctionValue               = "the given value is not function"
	ErrorCannotRegisterModule                = "failed to register module"
	ErrorTypeAlreadyRegistered               = "type already registered"
	ErrorServiceAlreadyLinkedWithAnotherList = "service already linked with another list"
	ErrorFailedToRegisterType                = "failed to register type"
)

var (

	// ErrRequiresFunctionValue indicates that the provided value is not a function.
	ErrRequiresFunctionValue = errors.New(ErrorRequiresFunctionValue)

	// ErrCannotRegisterModule indicates that the module registration process has failed.
	ErrCannotRegisterModule = errors.New(ErrorCannotRegisterModule)

	// ErrTypeAlreadyRegistered indicates that an attempt was made to register a type that is already registered.
	ErrTypeAlreadyRegistered = errors.New(ErrorTypeAlreadyRegistered)

	// ErrFailedToRegisterType indicates that the attempt to register a type has failed.
	ErrFailedToRegisterType = errors.New(ErrorFailedToRegisterType)
)

// RegistryError represents an error that gets returned for failing registry operations.
type RegistryError struct {
	ParsleyError
	serviceTypeName string
}

// _ ensures that RegistryError implements the ParsleyErrorWithServiceTypeName interface.
var _ ParsleyErrorWithServiceTypeName = &RegistryError{}

func (r *RegistryError) setServiceTypeName(name string) {
	r.serviceTypeName = name
}

// ServiceTypeName returns the service type name associated with the RegistryError.
func (r *RegistryError) ServiceTypeName() string {
	return r.serviceTypeName
}

// Error returns the message associated with the RegistryError.
func (r *RegistryError) Error() string {
	return r.Msg
}

// Format implements the fmt.Formatter interface.
func (r *RegistryError) Format(s fmt.State, verb rune) {
	formatError(r.Msg, r.serviceTypeName, r.cause, s, verb)
}

// MatchesServiceType checks if the service type name of the RegistryError matches the specified name.
func (r *RegistryError) MatchesServiceType(name string) bool {
	return r.serviceTypeName == name
}

// NewRegistryError creates a new RegistryError with the given message and initializers to modify the error.
func NewRegistryError(msg string, initializers ...ParsleyErrorFunc) error {
	err := &RegistryError{
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
