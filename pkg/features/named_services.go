package features

import (
	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/resolving"
	"github.com/matzefriedrich/parsley/pkg/types"
)

type namedService struct {
	name          string
	activatorFunc any
}

func (n namedService) ActivatorFunc() any {
	return n.activatorFunc
}

func (n namedService) Name() string {
	return n.name
}

func RegisterNamed[T any](registry types.ServiceRegistry, services ...registration.NamedServiceRegistrationFunc) error {

	registrationErrors := make([]error, 0)

	for _, service := range services {
		name, serviceActivatorFunc, _ := service()
		if len(name) == 0 || serviceActivatorFunc == nil {
			return types.NewRegistryError("invalid named service registration")
		}
		namedActivator := newNamedServiceFactory[T](name, serviceActivatorFunc)
		err := registration.RegisterInstance(registry, namedActivator)
		if err != nil {
			registrationErrors = append(registrationErrors, err)
		}
	}

	nameServiceResolver := resolving.CreateNamedServiceResolverActivatorFunc[T]()
	err := registration.RegisterTransient(registry, nameServiceResolver)
	if err != nil {
		registrationErrors = append(registrationErrors, err)
	}

	if len(registrationErrors) > 0 {
		return types.NewRegistryError("failed to register named services", types.WithAggregatedCause(registrationErrors...))
	}

	return nil
}

func newNamedServiceFactory[T any](name string, activatorFunc any) func() types.NamedService[T] {
	return func() types.NamedService[T] {
		return &namedService{
			name:          name,
			activatorFunc: activatorFunc,
		}
	}
}
