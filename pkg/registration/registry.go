package registration

import (
	"github.com/matzefriedrich/parsley/internal/core"
	"github.com/matzefriedrich/parsley/pkg/types"
	"reflect"
)

type serviceRegistry struct {
	identifierSource core.ServiceIdSequence
	registrations    map[reflect.Type]types.ServiceRegistrationList
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

func (s *serviceRegistry) addOrUpdateServiceRegistrationListFor(serviceType reflect.Type) types.ServiceRegistrationList {
	list, exists := s.registrations[serviceType]
	if exists {
		return list
	}
	list = NewServiceRegistrationList(s.identifierSource)
	s.registrations[serviceType] = list
	return list
}

func (s *serviceRegistry) Register(activatorFunc any, lifetimeScope types.LifetimeScope) error {

	registration, err := CreateServiceRegistration(activatorFunc, lifetimeScope)
	if err != nil {
		return err
	}

	serviceType := registration.ServiceType()
	list := s.addOrUpdateServiceRegistrationListFor(serviceType)
	addRegistrationErr := list.AddRegistration(registration)
	if addRegistrationErr != nil {
		return types.NewRegistryError("failed to register type", types.WithCause(addRegistrationErr))
	}

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

func (s *serviceRegistry) TryGetServiceRegistrations(serviceType reflect.Type) (types.ServiceRegistrationList, bool) {
	if s.IsRegistered(serviceType) == false {
		return nil, false
	}
	list, found := s.registrations[serviceType]
	if found && list.IsEmpty() == false {
		return list, true
	}
	return nil, false
}

func (s *serviceRegistry) TryGetSingleServiceRegistration(serviceType reflect.Type) (types.ServiceRegistration, bool) {
	list, found := s.TryGetServiceRegistrations(serviceType)
	if found && list.IsEmpty() == false {
		registrations := list.Registrations()
		const exactlyOne = 1
		if len(registrations) == exactlyOne {
			return registrations[0], true
		}
	}
	return nil, false
}

func NewServiceRegistry() types.ServiceRegistry {
	registrations := make(map[reflect.Type]types.ServiceRegistrationList)
	return &serviceRegistry{
		identifierSource: core.NewServiceId(0),
		registrations:    registrations,
	}
}

func (s *serviceRegistry) CreateLinkedRegistry() types.ServiceRegistry {
	registrations := make(map[reflect.Type]types.ServiceRegistrationList)
	return &serviceRegistry{
		identifierSource: s.identifierSource,
		registrations:    registrations,
	}
}

func (s *serviceRegistry) CreateScope() types.ServiceRegistry {
	registrations := make(map[reflect.Type]types.ServiceRegistrationList)
	for serviceType, registration := range s.registrations {
		registrations[serviceType] = registration
	}
	return &serviceRegistry{
		identifierSource: s.identifierSource,
		registrations:    registrations,
	}
}

var _ types.ServiceRegistry = &serviceRegistry{}
var _ types.ServiceRegistryAccessor = &serviceRegistry{}
