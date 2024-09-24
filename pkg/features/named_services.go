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

// ActivatorFunc retrieves the activator function associated with the named service.
func (n namedService) ActivatorFunc() any {
	return n.activatorFunc
}

// Name returns the name associated with the service instance. Used for identification purposes.
func (n namedService) Name() string {
	return n.name
}

// RegisterNamed registers named services with their respective activator functions and lifetime scopes.
// It supports dependency injection by associating names with service instances.
func RegisterNamed[T any](registry types.ServiceRegistry, services ...registration.NamedServiceRegistrationFunc) error {

	registrationErrors := make([]error, 0)

	for _, service := range services {
		name, serviceActivatorFunc, scope := service()
		if len(name) == 0 || serviceActivatorFunc == nil {
			return types.NewRegistryError("invalid named service registration")
		}
		_ = registry.Register(serviceActivatorFunc, scope)
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
