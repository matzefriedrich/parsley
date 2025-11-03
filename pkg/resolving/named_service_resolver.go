package resolving

import (
	"context"

	"github.com/matzefriedrich/parsley/pkg/types"
)

// NamedServiceResolverActivatorFunc defines a function for resolving named services.
type NamedServiceResolverActivatorFunc[T any] func(types.Resolver) func(string) (T, error)

// CreateNamedServiceResolverActivatorFunc creates a NamedServiceResolverActivatorFunc for resolving named services.
func CreateNamedServiceResolverActivatorFunc[T any](ctx context.Context) NamedServiceResolverActivatorFunc[T] {
	return func(resolver types.Resolver) func(string) (T, error) {
		var nilInstance T
		requiredServices, _ := ResolveRequiredServices[func() types.NamedService[T]](ctx, resolver)
		return func(name string) (T, error) {
			for _, service := range requiredServices {
				s := service()
				if s.Name() == name {
					return Activate[T](ctx, resolver, s.ActivatorFunc())
				}
			}
			return nilInstance, types.NewResolverError("failed to resolve named service")
		}
	}
}
