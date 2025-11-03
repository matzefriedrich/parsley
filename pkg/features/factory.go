package features

import (
	"context"

	"github.com/matzefriedrich/parsley/pkg/resolving"
	"github.com/matzefriedrich/parsley/pkg/types"
)

// FactoryFunc represents a function that creates an instance of type T using a context for dependency resolution.
type FactoryFunc[T any] func(ctx context.Context) (T, error)

// RegisterFactory registers a factory function for resolving instances of a specified type with a given lifetime scope.
func RegisterFactory[T any](registry types.ServiceRegistry, scope types.LifetimeScope) error {
	f := func(resolver types.Resolver) FactoryFunc[T] {
		return resolverFunc[T](resolver)
	}
	return registry.Register(f, scope)
}

// resolverFunc creates and returns a FactoryFunc that resolves an instance of type T using the provided resolver.
func resolverFunc[T any](resolver types.Resolver) FactoryFunc[T] {
	return func(ctx context.Context) (T, error) {
		return resolving.ResolveRequiredService[T](ctx, resolver)
	}
}
