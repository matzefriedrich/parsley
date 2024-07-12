package pkg

import (
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

func WithInstance[T any](instance T) types.ResolverOptionsFunc {
	return func(registry types.ServiceRegistry) error {
		err := RegisterInstance[T](registry, instance)
		if err != nil {
			return types.NewRegistryError(types.ErrorCannotRegisterTypeWithResolverOptions, types.WithCause(err))
		}
		return nil
	}
}
