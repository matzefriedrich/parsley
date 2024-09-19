package resolving

import (
	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/types"
)

func applyResolverOptions(registry types.ServiceRegistry, options ...types.ResolverOptionsFunc) error {
	for _, option := range options {
		err := option(registry)
		if err != nil {
			return err
		}
	}
	return nil
}

// WithInstance Creates a ResolverOptionsFunc that registers a specific instance of a type T with a service registry to be resolved as a singleton.
func WithInstance[T any](instance T) types.ResolverOptionsFunc {
	return func(registry types.ServiceRegistry) error {
		err := registration.RegisterInstance[T](registry, instance)
		if err != nil {
			return types.NewRegistryError(types.ErrorCannotRegisterTypeWithResolverOptions, types.WithCause(err))
		}
		return nil
	}
}
