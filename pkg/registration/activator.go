package registration

import (
	"github.com/matzefriedrich/parsley/internal"
	"github.com/matzefriedrich/parsley/pkg/types"
	"reflect"
)

// CreateServiceActivatorFrom creates a service activator function for a given instance of type T.
func CreateServiceActivatorFrom[T any](instance T) (func() T, error) {
	if internal.IsNil(instance) {
		return nil, types.NewRegistryError(types.ErrorInstanceCannotBeNil)
	}
	t := reflect.TypeOf((*T)(nil)).Elem()
	switch t.Kind() {
	case reflect.Func:
	case reflect.Interface:
	case reflect.Pointer:
	case reflect.Struct:
	default:
		return nil, types.NewRegistryError(types.ErrorActivatorFunctionInvalidReturnType)
	}
	instanceFunc := func() T {
		return instance
	}
	return instanceFunc, nil
}

// RegisterInstance registers an instance of type T. A registered instance behaves like a service registration with a singleton lifetime scope.
func RegisterInstance[T any](registry types.ServiceRegistry, instance T) error {
	instanceFunc, err := CreateServiceActivatorFrom[T](instance)
	if err != nil {
		return err
	}
	return registry.Register(instanceFunc, types.LifetimeSingleton)
}
