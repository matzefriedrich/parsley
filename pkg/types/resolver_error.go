package types

import "errors"

const (
	ErrorServiceTypeNotRegistered              = "service type is not registered"
	ErrorRequiredServiceNotRegistered          = "required service type is not registered"
	ErrorCannotResolveService                  = "cannot resolve service"
	ErrorActivatorFunctionInvalidReturnType    = "activator function has an invalid return type"
	ErrorCircularDependencyDetected            = "circular dependency detected"
	ErrorCannotBuildDependencyGraph            = "failed to build dependency graph"
	ErrorInstanceCannotBeNil                   = "instance cannot be nil"
	ErrorServiceTypeMustBeInterface            = "service type must be an interface"
	ErrorCannotRegisterTypeWithResolverOptions = "cannot register type with resolver options"
)

var (
	ErrServiceTypeNotRegistered              = errors.New(ErrorServiceTypeNotRegistered)
	ErrActivatorFunctionInvalidReturnType    = errors.New(ErrorCannotResolveService)
	ErrCannotBuildDependencyGraph            = errors.New(ErrorCannotBuildDependencyGraph)
	ErrCircularDependencyDetected            = errors.New(ErrorCircularDependencyDetected)
	ErrInstanceCannotBeNil                   = errors.New(ErrorInstanceCannotBeNil)
	ErrServiceTypeMustBeInterface            = errors.New(ErrorServiceTypeMustBeInterface)
	ErrCannotRegisterTypeWithResolverOptions = errors.New(ErrorCannotRegisterTypeWithResolverOptions)
)

type ResolverError struct {
	ParsleyError
	serviceTypeName string
}

var _ ParsleyErrorWithServiceTypeName = &ResolverError{}

func (r *ResolverError) ServiceTypeName(name string) {
	r.serviceTypeName = name
}

func NewResolverError(msg string, initializers ...ParsleyErrorFunc) error {
	err := &ResolverError{
		ParsleyError: ParsleyError{
			msg: msg,
		},
	}
	for _, f := range initializers {
		f(&err.ParsleyError)
		f(err)
	}

	return err
}
