package types

import "errors"

const (
	ErrorServiceTypeNotRegistered                = "service type is not registered"
	ErrorRequiredServiceNotRegistered            = "required service type is not registered"
	ErrorCannotResolveService                    = "cannot resolve service"
	ErrorActivatorFunctionsMustReturnAnInterface = "activator functions must return an interfaces"
	ErrorCircularDependencyDetected              = "circular dependency detected"
	ErrorCannotBuildDependencyGraph              = "failed to build dependency graph"
	ErrorInstanceCannotBeNil                     = "instance cannot be nil"
	ErrorServiceTypeMustBeInterface              = "service type must be an interface"
)

var (
	ErrServiceTypeNotRegistered                = errors.New(ErrorServiceTypeNotRegistered)
	ErrActivatorFunctionsMustReturnAnInterface = errors.New(ErrorCannotResolveService)
	ErrCannotBuildDependencyGraph              = errors.New(ErrorCannotBuildDependencyGraph)
	ErrCircularDependencyDetected              = errors.New(ErrorCircularDependencyDetected)
	ErrInstanceCannotBeNil                     = errors.New(ErrorInstanceCannotBeNil)
	ErrServiceTypeMustBeInterface              = errors.New(ErrorServiceTypeMustBeInterface)
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
