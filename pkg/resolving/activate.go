package resolving

import (
	"context"
	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/types"
)

// Activate attempts to create and return an instance of the requested type using the provided resolver.
// This method can be used to instantiate service objects of unregistered types. The specified activator function can
// have parameters to demand service instances for registered service types.
// See https://matzefriedrich.github.io/parsley-docs/resolving/resolve-live-services/ for further information.
func Activate[T any](resolver types.Resolver, ctx context.Context, activatorFunc any, options ...types.ResolverOptionsFunc) (T, error) {

	var nilInstance T

	lifetimeScope := types.LifetimeTransient
	serviceRegistration, registrationErr := registration.CreateServiceRegistration(activatorFunc, lifetimeScope)
	if registrationErr != nil {
		return nilInstance, types.NewResolverError(types.ErrorCannotCreateInstanceOfUnregisteredType, types.WithCause(registrationErr))
	}

	serviceType := serviceRegistration.ServiceType()
	resolveActivatorFuncOption := func(registry types.ServiceRegistry) error {
		return registry.Register(activatorFunc, lifetimeScope)
	}

	services, err := resolver.ResolveWithOptions(ctx, serviceType, resolveActivatorFuncOption)
	if err != nil {
		return nilInstance, err
	}

	if len(services) == 1 {
		compatible, ok := services[0].(T)
		if ok {
			return compatible, nil
		}
	} else if len(services) > 1 {
		return nilInstance, types.NewResolverError(types.ErrorAmbiguousServiceInstancesResolved)
	}

	return nilInstance, types.NewResolverError(types.ErrorCannotResolveService)
}
