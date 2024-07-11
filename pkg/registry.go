package pkg

import (
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

func RegisterInstance[T any](registry types.ServiceRegistry, instance T) error {
	instanceFunc, err := CreateServiceActivatorFrom[T](instance)
	if err != nil {
		return err
	}
	return registry.Register(instanceFunc, types.LifetimeSingleton)
}

func (s *serviceRegistry) Register(activatorFunc any, lifetimeScope types.LifetimeScope) error {

	registration, err := CreateServiceRegistration(activatorFunc, lifetimeScope)
	if err != nil {
		return err
	}

	id := s.identifierSource.Next()
	setupErr := registration.SetId(id)
	if setupErr != nil {
		return types.NewRegistryError("failed to set up type registration", types.WithCause(setupErr))
	}

	serviceType := registration.ServiceType()
	s.registrations[serviceType] = registration

	return nil
}

func (s *serviceRegistry) RegisterModule(modules ...types.ModuleFunc) error {
	for _, m := range modules {
		err := m(s)
		if err != nil {
			return types.NewRegistryError(types.ErrorCannotRegisterModule, types.WithCause(err))
		}
	}
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

func (s *serviceRegistry) CreateLinkedRegistry() types.ServiceRegistry {
	registrations := make(map[reflect.Type]types.ServiceRegistration)
	return &serviceRegistry{
		identifierSource: s.identifierSource,
		registrations:    registrations,
	}
}

func (s *serviceRegistry) CreateScope() types.ServiceRegistry {
	registrations := make(map[reflect.Type]types.ServiceRegistration)
	for serviceType, registration := range s.registrations {
		registrations[serviceType] = registration
	}
	return &serviceRegistry{
		identifierSource: s.identifierSource,
		registrations:    registrations,
	}
}

func (s *serviceRegistry) BuildResolver() types.Resolver {
	r := NewResolver(s)
	_ = RegisterInstance(s, r)
	return r
}

var _ types.ServiceRegistry = &serviceRegistry{}
var _ types.ServiceRegistryAccessor = &serviceRegistry{}
