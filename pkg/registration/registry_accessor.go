package registration

import (
	"github.com/matzefriedrich/parsley/pkg/types"
	"reflect"
)

type multiRegistryAccessor struct {
	registries []types.ServiceRegistryAccessor
}

func (m *multiRegistryAccessor) TryGetSingleServiceRegistration(serviceType reflect.Type) (types.ServiceRegistration, bool) {
	for _, registry := range m.registries {
		registration, ok := registry.TryGetSingleServiceRegistration(serviceType)
		if ok {
			return registration, ok
		}
	}
	return nil, false
}

func (m *multiRegistryAccessor) TryGetServiceRegistrations(serviceType reflect.Type) (types.ServiceRegistrationList, bool) {
	for _, registry := range m.registries {
		registration, ok := registry.TryGetServiceRegistrations(serviceType)
		if ok {
			return registration, ok
		}
	}
	return nil, false
}

var _ types.ServiceRegistryAccessor = &multiRegistryAccessor{}

func NewMultiRegistryAccessor(registries ...types.ServiceRegistryAccessor) types.ServiceRegistryAccessor {
	serviceRegistries := make([]types.ServiceRegistryAccessor, 0)
	serviceRegistries = append(serviceRegistries, registries...)
	return &multiRegistryAccessor{
		registries: serviceRegistries,
	}
}
