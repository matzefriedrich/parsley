package types

import (
	"context"
	"reflect"
)

type FunctionInfo interface {
	Name() string
	ParameterTypes() []reflect.Type
	ReturnType() reflect.Type
}

type ServiceRegistry interface {
	ServiceRegistryAccessor
	CreateLinkedRegistry() ServiceRegistry
	CreateScope() ServiceRegistry
	IsRegistered(serviceType reflect.Type) bool
	Register(activatorFunc any, scope LifetimeScope) error
	RegisterModule(modules ...ModuleFunc) error
}

type ModuleFunc func(registry ServiceRegistry) error

type ServiceRegistryAccessor interface {
	TryGetServiceRegistrations(serviceType reflect.Type) (ServiceRegistrationList, bool)
	TryGetSingleServiceRegistration(serviceType reflect.Type) (ServiceRegistration, bool)
}

type ServiceRegistration interface {
	Id() uint64
	InvokeActivator(params ...interface{}) (interface{}, error)
	IsSame(other ServiceRegistration) bool
	LifetimeScope() LifetimeScope
	RequiredServiceTypes() []reflect.Type
	ServiceType() reflect.Type
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

type RegistrationConfigurationFunc func(r ServiceRegistration)

type ResolverOptionsFunc func(registry ServiceRegistry) error

type Resolver interface {
	Resolve(ctx context.Context, serviceType reflect.Type) ([]interface{}, error)
	ResolveWithOptions(ctx context.Context, serviceType reflect.Type, options ...ResolverOptionsFunc) ([]interface{}, error)
}

type DependencyInfo interface {
	AddRequiredServiceInfo(child DependencyInfo)
	CreateInstance() (interface{}, error)
	Consumer() DependencyInfo
	HasInstance() bool
	Instance() interface{}
	Registration() ServiceRegistration
	RequiredServiceTypes() []reflect.Type
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
