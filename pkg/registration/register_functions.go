package registration

import "github.com/matzefriedrich/parsley/pkg/types"

// SupportsRegisterActivatorFunc allows the registration of activator functions with different lifetime scopes.
type SupportsRegisterActivatorFunc interface {
	Register(activatorFunc any, scope types.LifetimeScope) error
}

// SupportsRegisterActivatorFuncWithOptions allows the registration of activator functions with lifetime scope and lifecycle options.
type SupportsRegisterActivatorFuncWithOptions interface {
	RegisterWithOptions(activatorFunc any, scope types.LifetimeScope, options ...types.LifecycleOption) error
}

// RegisterTransientWithOptions registers services with a transient lifetime and the given lifecycle options.
func RegisterTransientWithOptions(registry SupportsRegisterActivatorFuncWithOptions, activatorFunc any, options ...types.LifecycleOption) error {
	return registry.RegisterWithOptions(activatorFunc, types.LifetimeTransient, options...)
}

// RegisterScopedWithOptions registers services with a scoped lifetime and the given lifecycle options.
func RegisterScopedWithOptions(registry SupportsRegisterActivatorFuncWithOptions, activatorFunc any, options ...types.LifecycleOption) error {
	return registry.RegisterWithOptions(activatorFunc, types.LifetimeScoped, options...)
}

// RegisterSingletonWithOptions registers services with a singleton lifetime and the given lifecycle options.
func RegisterSingletonWithOptions(registry SupportsRegisterActivatorFuncWithOptions, activatorFunc any, options ...types.LifecycleOption) error {
	return registry.RegisterWithOptions(activatorFunc, types.LifetimeSingleton, options...)
}

// RegisterTransient registers services with a transient lifetime in the provided service registry.
func RegisterTransient(registry SupportsRegisterActivatorFunc, activatorFunc ...any) error {
	for _, a := range activatorFunc {
		err := registry.Register(a, types.LifetimeTransient)
		if err != nil {
			return err
		}
	}
	return nil
}

// RegisterScoped registers services with a scoped lifetime in the provided service registry.
func RegisterScoped(registry SupportsRegisterActivatorFunc, activatorFunc ...any) error {
	for _, a := range activatorFunc {
		err := registry.Register(a, types.LifetimeScoped)
		if err != nil {
			return err
		}
	}
	return nil
}

// RegisterSingleton registers services with a singleton lifetime in the provided service registry.
func RegisterSingleton(registry SupportsRegisterActivatorFunc, activatorFunc ...any) error {
	for _, a := range activatorFunc {
		err := registry.Register(a, types.LifetimeSingleton)
		if err != nil {
			return err
		}
	}
	return nil
}
