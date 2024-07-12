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
	switch t.Kind() {
	case reflect.Func:
	case reflect.Interface:
	default:
		return nil, types.NewRegistryError(types.ErrorActivatorFunctionInvalidReturnType)
	}
	instanceFunc := func() T {
		return instance
	}
	return instanceFunc, nil
}
