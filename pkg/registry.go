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

func (s *serviceRegistry) Register(activatorFunc any, configuration ...types.ServiceConfigurationFunc) error {

	value := reflect.ValueOf(activatorFunc)

	info, err := reflectFunctionInfoFrom(value)
	if err != nil {
		return types.NewRegistryError(types.ErrorRequiresFunctionValue, types.WithCause(err))
	}

	serviceType := info.ReturnType()
	requiredTypes := info.ParameterTypes()

	registration := newServiceRegistration(serviceType, value, requiredTypes...)
	for _, config := range configuration {
		config(registration)
	}

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
