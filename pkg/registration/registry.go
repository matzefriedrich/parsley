package registration

import (
	"github.com/matzefriedrich/parsley/internal/core"
	"github.com/matzefriedrich/parsley/pkg/types"
)

type serviceRegistry struct {
	identifierSource core.ServiceIdSequence
	registrations    map[types.ServiceKey]types.ServiceRegistrationList
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

func (s *serviceRegistry) addOrUpdateServiceRegistrationListFor(serviceType types.ServiceType) types.ServiceRegistrationList {
	list, exists := s.registrations[serviceType.LookupKey()]
	if exists {
		return list
	}
	list = NewServiceRegistrationList(s.identifierSource)
	s.registrations[serviceType.LookupKey()] = list
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
		return types.NewRegistryError(types.ErrorFailedToRegisterType, types.WithCause(addRegistrationErr))
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

func (s *serviceRegistry) IsRegistered(serviceType types.ServiceType) bool {
	_, found := s.registrations[serviceType.LookupKey()]
	return found
}

func (s *serviceRegistry) TryGetServiceRegistrations(serviceType types.ServiceType) (types.ServiceRegistrationList, bool) {
	if s.IsRegistered(serviceType) == false {
		return nil, false
	}
	list, found := s.registrations[serviceType.LookupKey()]
	if found && list.IsEmpty() == false {
		return list, true
	}
	return nil, false
}

func (s *serviceRegistry) TryGetSingleServiceRegistration(serviceType types.ServiceType) (types.ServiceRegistration, bool) {
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
	registrations := make(map[types.ServiceKey]types.ServiceRegistrationList)
	return &serviceRegistry{
		identifierSource: core.NewServiceId(0),
		registrations:    registrations,
	}
}

func (s *serviceRegistry) CreateLinkedRegistry() types.ServiceRegistry {
	registrations := make(map[types.ServiceKey]types.ServiceRegistrationList)
	return &serviceRegistry{
		identifierSource: s.identifierSource,
		registrations:    registrations,
	}
}

func (s *serviceRegistry) CreateScope() types.ServiceRegistry {
	registrations := make(map[types.ServiceKey]types.ServiceRegistrationList)
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
