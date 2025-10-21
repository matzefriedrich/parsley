package types

import (
	"context"
	"fmt"
	"reflect"
)

// FunctionInfo Stores information about a service activator function. This interface supports the internal infrastructure.
type FunctionInfo interface {
	fmt.Stringer
	Name() string
	Parameters() []FunctionParameterInfo
	ReturnType() ServiceType
	ParameterTypes() []ServiceType
}

type FunctionParameterInfo interface {
	fmt.Stringer
	Type() ServiceType
}

// ServiceKey represents a unique key for identifying services in the service registry.
type ServiceKey struct {
	value string
}

// String Gets the value of the current ServiceKey instance.
func (s ServiceKey) String() string {
	return s.value
}

// NewServiceKey creates a new ServiceKey with the given value.
func NewServiceKey(value string) ServiceKey {
	return ServiceKey{value: value}
}

// ServiceType represents a service type.
type ServiceType interface {

	// Name returns the name of the service type.
	Name() string

	// PackagePath returns the package path of the service type.
	PackagePath() string

	// ReflectedType returns the underlying reflect.Type representation of the service type.
	ReflectedType() reflect.Type

	// LookupKey retrieves the ServiceKey associated with the service type.
	LookupKey() ServiceKey
}

// ServiceRegistry provides methods to map service types to activator functions. The service registration organizes and stores the metadata required by the service resolver.
type ServiceRegistry interface {
	ServiceRegistryAccessor

	// CreateLinkedRegistry creates and returns a new ServiceRegistry instance linked to the current registry. A linked service registry is an empty service registry.
	CreateLinkedRegistry() ServiceRegistry

	// CreateScope creates and returns a scoped ServiceRegistry instance which inherits all service registrations from the current ServiceRegistry instance.
	CreateScope() ServiceRegistry

	// GetServiceRegistrations retrieves all service registrations.
	GetServiceRegistrations() ([]ServiceRegistration, error)

	// IsRegistered checks if a service of the specified ServiceType is registered in the service registry.
	IsRegistered(serviceType ServiceType) bool

	// Register registers a service with its activator function and defines its lifetime scope with the service registry.
	Register(activatorFunc any, scope LifetimeScope) error

	// RegisterModule registers one or more modules, encapsulated as ModuleFunc, with the service registry. A module is a logical unit of service registrations.
	RegisterModule(modules ...ModuleFunc) error
}

// ModuleFunc defines a function used to register services with the given service registry.
type ModuleFunc func(registry ServiceRegistry) error

// ServiceRegistryAccessor provides methods to access and retrieve service registrations from the registry.
type ServiceRegistryAccessor interface {

	// TryGetServiceRegistrations attempts to retrieve all service registrations for the given service type.
	// Returns the service registration list and true if found, otherwise returns false.
	TryGetServiceRegistrations(serviceType ServiceType) (ServiceRegistrationList, bool)

	// TryGetSingleServiceRegistration attempts to retrieve a single service registration for the given service type.
	// Returns the service registration and true if found, otherwise returns false.
	TryGetSingleServiceRegistration(serviceType ServiceType) (ServiceRegistration, bool)
}

// ServiceRegistration represents a service registrations.
type ServiceRegistration interface {

	// Id Returns the unique identifier of the service registration.
	Id() uint64

	// InvokeActivator calls the activator function with the provided parameters and returns the resulting instance and any error.
	InvokeActivator(params ...interface{}) (interface{}, error)

	// IsSame checks if the provided ServiceRegistration equals the current ServiceRegistration.
	IsSame(other ServiceRegistration) bool

	// LifetimeScope returns the LifetimeScope associated with the service registration.
	LifetimeScope() LifetimeScope

	// RequiredServiceTypes returns a slice of ServiceType, containing all service types required by the service registration.
	RequiredServiceTypes() []ServiceType

	// ServiceType retrieves the type of the service being registered.
	ServiceType() ServiceType
}

// ServiceRegistrationList provides functionality to manage a list of service registrations. This interface supports internal infrastructure services.
type ServiceRegistrationList interface {

	// AddRegistration adds a new service registration to the list.
	AddRegistration(registration ServiceRegistrationSetup) error

	// Id returns the unique identifier of the service registration list.
	Id() uint64

	// Registrations returns a slice of ServiceRegistration, containing all registrations in the list.
	Registrations() []ServiceRegistration

	// IsEmpty checks if the service registration list contains any registrations.
	// It returns true if the list is empty, otherwise false.
	IsEmpty() bool
}

// ServiceRegistrationSetup extends ServiceRegistration and supports internal infrastructure services.
type ServiceRegistrationSetup interface {
	ServiceRegistration

	// SetId sets the unique identifier for the service registration. This method supports internal infrastructure and is not intended to be used by your code.
	SetId(id uint64) error
}

// NamedService is a generic interface defining a service with a name and an activator function.
type NamedService[T any] interface {
	Name() string
	ActivatorFunc() any
}

// ResolverOptionsFunc represents a function that configures a service registry used by the resolver.
type ResolverOptionsFunc func(registry ServiceRegistry) error

// Resolver provides methods to resolve registered services based on types.
type Resolver interface {

	// Resolve attempts to resolve all registered services of the specified ServiceType.
	Resolve(ctx context.Context, serviceType ServiceType) ([]interface{}, error)

	// ResolveWithOptions resolves services of the specified type using additional options and returns a list of resolved services or an error.
	ResolveWithOptions(ctx context.Context, serviceType ServiceType, options ...ResolverOptionsFunc) ([]interface{}, error)
}

// DependencyInfo provides functionality to manage dependency information.
type DependencyInfo interface {
	// AddRequiredServiceInfo adds a child dependency to the current dependency info.
	AddRequiredServiceInfo(child DependencyInfo)

	// CreateInstance creates an instance of the service associated with this dependency info.
	CreateInstance() (interface{}, error)

	// Consumer returns the parent dependency for the current dependency info.
	Consumer() DependencyInfo

	// HasInstance checks if an instance has already been created for the dependency represented by the current DependencyInfo object.
	HasInstance() bool

	// Instance retrieves the created instance of the service associated with this dependency info.
	Instance() interface{}

	// Registration gets the service registration of the current dependency info.
	Registration() ServiceRegistration

	// RequiredServiceTypes gets the service types required by this dependency info.
	RequiredServiceTypes() []ServiceType

	// RequiredServices retrieves the instances of services required by this dependency info.
	RequiredServices() ([]interface{}, error)

	// ServiceTypeName gets the name of the service type associated with this dependency info.
	ServiceTypeName() string

	// SetInstance sets the instance for the current dependency info.
	SetInstance(instance interface{}) error
}

// LifetimeScope represents the duration for which a service or object instance is retained.
type LifetimeScope uint

func (l LifetimeScope) String() string {
	switch l {
	case LifetimeTransient:
		return "transient"
	case LifetimeScoped:
		return "scoped"
	case LifetimeSingleton:
		return "singleton"
	}
	return ""
}

const (

	// LifetimeTransient represents a transient lifetime where a new instance is created each time it is requested.
	LifetimeTransient LifetimeScope = iota

	// LifetimeScoped represents a scoped lifetime where a single instance is created per scope.
	LifetimeScoped

	// LifetimeSingleton represents a single instance scope that persists for the lifetime of the application.
	LifetimeSingleton
)
