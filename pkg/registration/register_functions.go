package registration

import "github.com/matzefriedrich/parsley/pkg/types"

// RegisterTransient registers services with a transient lifetime in the provided service registry.
// See https://matzefriedrich.github.io/parsley-docs/registration/register-constructor-functions/ for further information.
func RegisterTransient(registry types.ServiceRegistry, activatorFunc ...any) error {
	for _, a := range activatorFunc {
		err := registry.Register(a, types.LifetimeTransient)
		if err != nil {
			return err
		}
	}
	return nil
}

// RegisterScoped registers services with a scoped lifetime in the provided service registry.
func RegisterScoped(registry types.ServiceRegistry, activatorFunc ...any) error {
	for _, a := range activatorFunc {
		err := registry.Register(a, types.LifetimeScoped)
		if err != nil {
			return err
		}
	}
	return nil
}

// RegisterSingleton registers services with a singleton lifetime in the provided service registry.
func RegisterSingleton(registry types.ServiceRegistry, activatorFunc ...any) error {
	for _, a := range activatorFunc {
		err := registry.Register(a, types.LifetimeSingleton)
		if err != nil {
			return err
		}
	}
	return nil
}
