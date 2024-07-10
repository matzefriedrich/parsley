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
	BuildResolver() Resolver
	IsRegistered(serviceType reflect.Type) bool
	Register(activatorFunc any, scope LifetimeScope) error
}

type ServiceRegistryAccessor interface {
	ServiceRegistry
	TryGetServiceRegistration(serviceType reflect.Type) (ServiceRegistration, bool)
}

type ServiceRegistration interface {
	Id() uint64
	InvokeActivator(params ...interface{}) (interface{}, error)
	RequiredServiceTypes() []reflect.Type
	ServiceType() reflect.Type
	LifetimeScope() LifetimeScope
}

type RegistrationConfigurationFunc func(r ServiceRegistration)

type Resolver interface {
	Resolve(ctx context.Context, serviceType reflect.Type) (interface{}, error)
}

type DependencyInfo interface {
	AddRequiredServiceInfo(child DependencyInfo)
	CreateInstance() (interface{}, error)
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
