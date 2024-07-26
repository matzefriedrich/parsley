package types

import (
	"reflect"
)

type serviceType struct {
	reflectedType reflect.Type
	name          string
}

func (s serviceType) ReflectedType() reflect.Type {
	return s.reflectedType
}

func (s serviceType) Name() string {
	return s.name
}

func MakeServiceType[T any]() ServiceType {
	elem := reflect.TypeOf(new(T)).Elem()
	return &serviceType{
		reflectedType: elem,
		name:          elem.String(),
	}
}

func ServiceTypeFrom(t reflect.Type) ServiceType {
	name := ""
	switch t.Kind() {
	case reflect.Ptr:
		name = t.Elem().String()
	case reflect.Interface:
		name = t.Name()
	case reflect.Func:
		name = t.String()
	default:
		panic("unsupported type: " + t.String())
	}
	return &serviceType{
		reflectedType: t,
		name:          name,
	}
}
