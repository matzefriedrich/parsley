package registration

import (
	"github.com/matzefriedrich/parsley/internal/core"
	"github.com/matzefriedrich/parsley/pkg/types"
)

type serviceRegistry struct {
	identifierSource core.ServiceIdSequence
	registrations    map[types.ServiceKey]types.ServiceRegistrationList
}

// RegisterTransient adds a transient service registration with the provided activator function.
// See https://matzefriedrich.github.io/parsley-docs/registration/register-constructor-functions/ for further information.
func RegisterTransient(registry types.ServiceRegistry, activatorFunc any) error {
	return registry.Register(activatorFunc, types.LifetimeTransient)
}

// RegisterScoped adds a scoped service registration with the provided activator function.
func RegisterScoped(registry types.ServiceRegistry, activatorFunc any) error {
	return registry.Register(activatorFunc, types.LifetimeScoped)
}

// RegisterSingleton adds a singleton service registration with the provided activator function.
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

// GetServiceRegistrations retrieves all registered services as a slice of ServiceRegistration.
func (s *serviceRegistry) GetServiceRegistrations() ([]types.ServiceRegistration, error) {
	registrations := make([]types.ServiceRegistration, 0)
	for _, list := range s.registrations {
		registrations = append(registrations, list.Registrations()...)
	}
	return registrations, nil
}

// Register adds a service registration with the provided activator function and lifetime scope.
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

// RegisterModule registers one or more modules with the service registry.
func (s *serviceRegistry) RegisterModule(modules ...types.ModuleFunc) error {
	for _, m := range modules {
		err := m(s)
		if err != nil {
			return types.NewRegistryError(types.ErrorCannotRegisterModule, types.WithCause(err))
		}
	}
	return nil
}

// IsRegistered checks if a service type is registered in the service registry.
func (s *serviceRegistry) IsRegistered(serviceType types.ServiceType) bool {
	_, found := s.registrations[serviceType.LookupKey()]
	return found
}

// TryGetServiceRegistrations Tries to find all available service registrations for the given service type.
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

// TryGetSingleServiceRegistration Tries to find a single service registration for the given service type.
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

// NewServiceRegistry creates a new types.ServiceRegistry instance.
func NewServiceRegistry() types.ServiceRegistry {
	registrations := make(map[types.ServiceKey]types.ServiceRegistrationList)
	return &serviceRegistry{
		identifierSource: core.NewServiceId(0),
		registrations:    registrations,
	}
}

// CreateLinkedRegistry creates and returns a new, empty ServiceRegistry instance linked to the current registry.
func (s *serviceRegistry) CreateLinkedRegistry() types.ServiceRegistry {
	registrations := make(map[types.ServiceKey]types.ServiceRegistrationList)
	return &serviceRegistry{
		identifierSource: s.identifierSource,
		registrations:    registrations,
	}
}

// CreateScope creates and returns a scoped types.ServiceRegistry instance that inherits all service registrations from the current registry.
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
