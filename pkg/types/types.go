package types

import (
	"context"
	"reflect"
)

type FunctionInfo interface {
	Name() string
	ReturnType() reflect.Type
	ParameterTypes() []reflect.Type
}

type ServiceRegistry interface {
	BuildResolver() Resolver
	Register(activatorFunc any, configuration ...ServiceConfigurationFunc) error
	IsRegistered(serviceType reflect.Type) bool
}

type ServiceRegistryAccessor interface {
	ServiceRegistry
	TryGetServiceRegistration(serviceType reflect.Type) (ServiceRegistration, bool)
}

type ServiceRegistration interface {
	Id() uint64
	ServiceType() reflect.Type
	RequiredServiceTypes() []reflect.Type
	InvokeActivator(params ...interface{}) (interface{}, error)
}

type ServiceConfigurationFunc func(r ServiceRegistration)

type Resolver interface {
	Resolve(ctx context.Context, serviceType reflect.Type) (interface{}, error)
}

type DependencyInfo interface {
	AddRequiredServiceInfo(child DependencyInfo)
	RequiredServiceTypes() []reflect.Type
	RequiredServices() ([]interface{}, error)
	CreateInstance() (interface{}, error)
	HasInstance() bool
	Instance() interface{}
	ServiceTypeName() string
}
