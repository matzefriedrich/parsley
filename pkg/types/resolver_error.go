package types

import "errors"

const (
	ErrorServiceTypeNotRegistered               = "service type is not registered"
	ErrorRequiredServiceNotRegistered           = "required service type is not registered"
	ErrorCannotResolveService                   = "cannot resolve service"
	ErrorAmbiguousServiceInstancesResolved      = "the resolve operation resulted in multiple service instances"
	ErrorActivatorFunctionInvalidReturnType     = "activator function has an invalid return type"
	ErrorCircularDependencyDetected             = "circular dependency detected"
	ErrorCannotBuildDependencyGraph             = "failed to build dependency graph"
	ErrorInstanceCannotBeNil                    = "instance cannot be nil"
	ErrorServiceTypeMustBeInterface             = "service type must be an interface"
	ErrorCannotRegisterTypeWithResolverOptions  = "cannot register type with resolver options"
	ErrorCannotCreateInstanceOfUnregisteredType = "failed to create instance of unregistered type"
)

var (
	ErrServiceTypeNotRegistered               = errors.New(ErrorServiceTypeNotRegistered)
	ErrActivatorFunctionInvalidReturnType     = errors.New(ErrorCannotResolveService)
	ErrCannotBuildDependencyGraph             = errors.New(ErrorCannotBuildDependencyGraph)
	ErrCircularDependencyDetected             = errors.New(ErrorCircularDependencyDetected)
	ErrInstanceCannotBeNil                    = errors.New(ErrorInstanceCannotBeNil)
	ErrServiceTypeMustBeInterface             = errors.New(ErrorServiceTypeMustBeInterface)
	ErrCannotRegisterTypeWithResolverOptions  = errors.New(ErrorCannotRegisterTypeWithResolverOptions)
	ErrCannotCreateInstanceOfUnregisteredType = errors.New(ErrorCannotCreateInstanceOfUnregisteredType)
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
			Msg: msg,
		},
	}
	for _, f := range initializers {
		f(&err.ParsleyError)
		f(err)
	}

	return err
}
