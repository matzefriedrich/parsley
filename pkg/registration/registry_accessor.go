package registration

import (
	"github.com/matzefriedrich/parsley/pkg/types"
	"reflect"
)

type multiRegistryAccessor struct {
	registries []types.ServiceRegistryAccessor
}

func (m *multiRegistryAccessor) TryGetServiceRegistration(serviceType reflect.Type) (types.ServiceRegistration, bool) {
	for _, registry := range m.registries {
		registration, ok := registry.TryGetServiceRegistration(serviceType)
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
