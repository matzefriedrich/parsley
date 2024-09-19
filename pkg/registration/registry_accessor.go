package registration

import (
	"github.com/matzefriedrich/parsley/pkg/types"
)

type multiRegistryAccessor struct {
	registries []types.ServiceRegistryAccessor
}

// TryGetSingleServiceRegistration attempts to find a single service registration for the given service type in multiple registries.
func (m *multiRegistryAccessor) TryGetSingleServiceRegistration(serviceType types.ServiceType) (types.ServiceRegistration, bool) {
	for _, registry := range m.registries {
		registration, ok := registry.TryGetSingleServiceRegistration(serviceType)
		if ok {
			return registration, ok
		}
	}
	return nil, false
}

// TryGetServiceRegistrations tries to retrieve all service registrations for the given service type from multiple registries.
func (m *multiRegistryAccessor) TryGetServiceRegistrations(serviceType types.ServiceType) (types.ServiceRegistrationList, bool) {
	for _, registry := range m.registries {
		registration, ok := registry.TryGetServiceRegistrations(serviceType)
		if ok {
			return registration, ok
		}
	}
	return nil, false
}

var _ types.ServiceRegistryAccessor = &multiRegistryAccessor{}

// NewMultiRegistryAccessor creates a new ServiceRegistryAccessor that aggregates multiple registries.
func NewMultiRegistryAccessor(registries ...types.ServiceRegistryAccessor) types.ServiceRegistryAccessor {
	serviceRegistries := make([]types.ServiceRegistryAccessor, 0)
	serviceRegistries = append(serviceRegistries, registries...)
	return &multiRegistryAccessor{
		registries: serviceRegistries,
	}
}
