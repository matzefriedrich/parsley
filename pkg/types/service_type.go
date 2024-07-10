package types

import "reflect"

func ServiceType[T any]() reflect.Type {
	return reflect.TypeOf(new(T)).Elem()
}
