package types

import (
	"fmt"
	"reflect"
)

type serviceType struct {
	reflectedType reflect.Type
	name          string
	packagePath   string
	list          bool
	lookupKey     ServiceKey
}

func (s serviceType) String() string {
	return fmt.Sprintf("Name: \"%s\", Package: \"%s\", List: %t)", s.name, s.packagePath, s.list)
}

var _ ServiceType = &serviceType{}

func (s serviceType) LookupKey() ServiceKey {
	return s.lookupKey
}

func (s serviceType) ReflectedType() reflect.Type {
	return s.reflectedType
}

func (s serviceType) Name() string {
	return s.name
}

func (s serviceType) PackagePath() string {
	return s.packagePath
}

func MakeServiceType[T any]() ServiceType {
	elem := reflect.TypeOf(new(T)).Elem()
	return ServiceTypeFrom(elem)
}

func ServiceTypeFrom(t reflect.Type) ServiceType {
	isList := false
	elemType := t
	switch t.Kind() {
	case reflect.Ptr:
		elemType = t.Elem()
	case reflect.Interface:
		break
	case reflect.Func:
		break
	case reflect.Slice:
		t = t.Elem()
		isList = true
	default:
		panic("unsupported type: " + t.String())
	}
	return newServiceType(t, elemType, isList)
}

func newServiceType(t reflect.Type, elemType reflect.Type, isList bool) ServiceType {
	packagePath := t.PkgPath()
	name := elemType.Name()
	key := fmt.Sprintf("%s.%s", packagePath, name)
	if isList {
		key = fmt.Sprintf("%s.%s[]", packagePath, name)
	}
	serviceKey := ServiceKey{value: key}
	return &serviceType{
		reflectedType: t,
		name:          name,
		packagePath:   packagePath,
		list:          isList,
		lookupKey:     serviceKey,
	}
}
