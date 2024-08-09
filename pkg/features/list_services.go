package features

import (
	"context"
	"github.com/matzefriedrich/parsley/pkg/resolving"
	"github.com/matzefriedrich/parsley/pkg/types"
)

func RegisterList[T any](registry types.ServiceRegistry) error {
	return registry.Register(func(resolver types.Resolver) []T {
		services, _ := resolving.ResolveRequiredServices[T](resolver, context.Background())
		return services
	}, types.LifetimeTransient)
}
