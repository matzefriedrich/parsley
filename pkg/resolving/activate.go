package resolving

import (
	"context"
	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/types"
)

func Activate[T any](resolver types.Resolver, ctx context.Context, activatorFunc any, options ...types.ResolverOptionsFunc) (T, error) {

	var nilInstance T

	lifetimeScope := types.LifetimeTransient
	registration, registrationErr := registration.CreateServiceRegistration(activatorFunc, lifetimeScope)
	if registrationErr != nil {
		return nilInstance, types.NewResolverError(types.ErrorCannotCreateInstanceOfUnregisteredType, types.WithCause(registrationErr))
	}

	serviceType := registration.ServiceType()
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
