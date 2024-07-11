package pkg

import (
	"github.com/matzefriedrich/parsley/internal"
	"github.com/matzefriedrich/parsley/pkg/types"
	"reflect"
)

func CreateServiceActivatorFrom[T any](instance T) (func() T, error) {
	if internal.IsNil(instance) {
		return nil, types.NewRegistryError(types.ErrorInstanceCannotBeNil)
	}
	t := reflect.TypeOf((*T)(nil)).Elem()
	if t.Kind() != reflect.Interface {
		return nil, types.NewRegistryError(types.ErrorActivatorFunctionsMustReturnAnInterface)
	}
	instanceFunc := func() T {
		return instance
	}
	return instanceFunc, nil
}
