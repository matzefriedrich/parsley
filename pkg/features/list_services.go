package features

import (
	"context"

	"github.com/matzefriedrich/parsley/pkg/resolving"
	"github.com/matzefriedrich/parsley/pkg/types"
)

// RegisterList registers a function that resolves and returns a list of services of type T with the specified registry.
func RegisterList[T any](ctx context.Context, registry types.ServiceRegistry) error {
	return registry.Register(func(resolver types.Resolver) []T {
		services, _ := resolving.ResolveRequiredServices[T](ctx, resolver)
		return services
	}, types.LifetimeTransient)
}
