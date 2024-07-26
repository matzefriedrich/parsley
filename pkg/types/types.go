package types

import (
	"context"
	"fmt"
	"reflect"
)

type FunctionInfo interface {
	fmt.Stringer
	Name() string
	ParameterTypes() []ServiceType
	ReturnType() ServiceType
}

type ServiceType interface {
	Name() string
	ReflectedType() reflect.Type
}

type ServiceRegistry interface {
	ServiceRegistryAccessor
	CreateLinkedRegistry() ServiceRegistry
	CreateScope() ServiceRegistry
	IsRegistered(serviceType ServiceType) bool
	Register(activatorFunc any, scope LifetimeScope) error
	RegisterModule(modules ...ModuleFunc) error
}

type ModuleFunc func(registry ServiceRegistry) error

type ServiceRegistryAccessor interface {
	TryGetServiceRegistrations(serviceType ServiceType) (ServiceRegistrationList, bool)
	TryGetSingleServiceRegistration(serviceType ServiceType) (ServiceRegistration, bool)
}

type ServiceRegistration interface {
	Id() uint64
	InvokeActivator(params ...interface{}) (interface{}, error)
	IsSame(other ServiceRegistration) bool
	LifetimeScope() LifetimeScope
	RequiredServiceTypes() []ServiceType
	ServiceType() ServiceType
}

type ServiceRegistrationList interface {
	AddRegistration(registration ServiceRegistrationSetup) error
	Id() uint64
	Registrations() []ServiceRegistration
	IsEmpty() bool
}

type ServiceRegistrationSetup interface {
	ServiceRegistration
	SetId(id uint64) error
}

type NamedService[T any] interface {
	Name() string
	ActivatorFunc() any
}

type RegistrationConfigurationFunc func(r ServiceRegistration)

type ResolverOptionsFunc func(registry ServiceRegistry) error

type Resolver interface {
	Resolve(ctx context.Context, serviceType ServiceType) ([]interface{}, error)
	ResolveWithOptions(ctx context.Context, serviceType ServiceType, options ...ResolverOptionsFunc) ([]interface{}, error)
}

type DependencyInfo interface {
	AddRequiredServiceInfo(child DependencyInfo)
	CreateInstance() (interface{}, error)
	Consumer() DependencyInfo
	HasInstance() bool
	Instance() interface{}
	Registration() ServiceRegistration
	RequiredServiceTypes() []ServiceType
	RequiredServices() ([]interface{}, error)
	ServiceTypeName() string
	SetInstance(instance interface{}) error
}

type LifetimeScope uint

const (
	LifetimeTransient LifetimeScope = iota
	LifetimeScoped
	LifetimeSingleton
)
