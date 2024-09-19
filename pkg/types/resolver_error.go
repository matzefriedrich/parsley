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

	// ErrServiceTypeNotRegistered is returned when attempting to resolve a service type that has not been registered.
	ErrServiceTypeNotRegistered = errors.New(ErrorServiceTypeNotRegistered)

	// ErrRequiredServiceNotRegistered is returned when a required service type is not registered.
	ErrRequiredServiceNotRegistered = errors.New(ErrorRequiredServiceNotRegistered)

	// ErrActivatorFunctionInvalidReturnType is returned when an activator function has an invalid return type.
	ErrActivatorFunctionInvalidReturnType = errors.New(ErrorCannotResolveService)

	// ErrCannotBuildDependencyGraph is returned when the resolver fails to build a dependency graph due to missing dependencies or other issues.
	ErrCannotBuildDependencyGraph = errors.New(ErrorCannotBuildDependencyGraph)

	// ErrCircularDependencyDetected is returned when a circular dependency is detected during the resolution process.
	ErrCircularDependencyDetected = errors.New(ErrorCircularDependencyDetected)

	// ErrInstanceCannotBeNil is returned when an instance provided is nil, but a non-nil value is required.
	ErrInstanceCannotBeNil = errors.New(ErrorInstanceCannotBeNil)

	// ErrServiceTypeMustBeInterface is returned when a service type is not an interface.
	ErrServiceTypeMustBeInterface = errors.New(ErrorServiceTypeMustBeInterface)

	// ErrCannotRegisterTypeWithResolverOptions is returned when the resolver failed to register a type via resolver options.
	ErrCannotRegisterTypeWithResolverOptions = errors.New(ErrorCannotRegisterTypeWithResolverOptions)

	// ErrCannotCreateInstanceOfUnregisteredType is returned when the resolver fails to instantiate a type that has not been registered.
	ErrCannotCreateInstanceOfUnregisteredType = errors.New(ErrorCannotCreateInstanceOfUnregisteredType)
)

// ResolverError represents an error that gets returned for failing service resolver operations.
type ResolverError struct {
	ParsleyError
	serviceTypeName string
}

// _ ensures that the ResolverError implements the ParsleyErrorWithServiceTypeName interface.
var _ ParsleyErrorWithServiceTypeName = &ResolverError{}

// ServiceTypeName sets the service type name for the ResolverError instance.
func (r *ResolverError) ServiceTypeName(name string) {
	r.serviceTypeName = name
}

// NewResolverError creates a new ResolverError with the provided message and applies optional ParsleyErrorFunc initializers.
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
