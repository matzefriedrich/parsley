package pkg

import (
	"errors"
	"github.com/matzefriedrich/parsley/internal"
	"github.com/matzefriedrich/parsley/pkg/types"
	"reflect"
)

type serviceRegistry struct {
	identifierSource internal.ServiceIdSequence
	registrations    map[reflect.Type]types.ServiceRegistration
}

func RegisterTransient(registry types.ServiceRegistry, activatorFunc any) error {
	return registry.Register(activatorFunc, types.LifetimeTransient)
}

func RegisterScoped(registry types.ServiceRegistry, activatorFunc any) error {
	return registry.Register(activatorFunc, types.LifetimeScoped)
}

func RegisterSingleton(registry types.ServiceRegistry, activatorFunc any) error {
	return registry.Register(activatorFunc, types.LifetimeSingleton)
}

func (s *serviceRegistry) Register(activatorFunc any, lifetimeScope types.LifetimeScope) error {

	value := reflect.ValueOf(activatorFunc)

	info, err := reflectFunctionInfoFrom(value)
	if err != nil {
		return types.NewRegistryError(types.ErrorRequiresFunctionValue, types.WithCause(err))
	}

	serviceType := info.ReturnType()
	if serviceType.Kind() != reflect.Interface {
		return errors.New(types.ErrorActivatorFunctionsMustReturnAnInterface)
	}

	requiredTypes := info.ParameterTypes()

	registration := newServiceRegistration(serviceType, lifetimeScope, value, requiredTypes...)

	registration.id = s.identifierSource.Next()
	s.registrations[serviceType] = registration

	return nil
}

func (s *serviceRegistry) IsRegistered(serviceType reflect.Type) bool {
	_, found := s.registrations[serviceType]
	return found
}

func (s *serviceRegistry) TryGetServiceRegistration(serviceType reflect.Type) (types.ServiceRegistration, bool) {
	if s.IsRegistered(serviceType) == false {
		return nil, false
	}
	registration, found := s.registrations[serviceType]
	if found {
		return registration, true
	}
	return nil, false
}

func NewServiceRegistry() types.ServiceRegistry {
	registrations := make(map[reflect.Type]types.ServiceRegistration)
	return &serviceRegistry{
		identifierSource: internal.NewServiceId(0),
		registrations:    registrations,
	}
}

func (s *serviceRegistry) BuildResolver() types.Resolver {
	return NewResolver(s)
}

var _ types.ServiceRegistry = &serviceRegistry{}
var _ types.ServiceRegistryAccessor = &serviceRegistry{}
