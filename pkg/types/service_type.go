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

// String returns a formatted string representation of the service type, including its name, and package path.
func (s serviceType) String() string {
	return fmt.Sprintf("Name: \"%s\", Package: \"%s\", List: %t)", s.name, s.packagePath, s.list)
}

var _ ServiceType = &serviceType{}

// LookupKey returns the lookup key associated with the service type.
func (s serviceType) LookupKey() ServiceKey {
	return s.lookupKey
}

// ReflectedType returns the reflected type of the service.
func (s serviceType) ReflectedType() reflect.Type {
	return s.reflectedType
}

// Name returns the name of the service type.
func (s serviceType) Name() string {
	return s.name
}

// PackagePath returns the package path of the service type.
func (s serviceType) PackagePath() string {
	return s.packagePath
}

// MakeServiceType creates a ServiceType instance for the specified generic type T.
func MakeServiceType[T any]() ServiceType {
	elem := reflect.TypeOf(new(T)).Elem()
	return ServiceTypeFrom(elem)
}

// ServiceTypeFrom creates a ServiceType from the given reflect.Type.
// Supports pointer, interface, function, slice, and struct types. The function panics, if t is of an unsupported kind is given.
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
	case reflect.Struct:
	default:
		panic("unsupported type: " + t.String())
	}
	return newServiceType(t, elemType, isList)
}

func newServiceType(t reflect.Type, elemType reflect.Type, isList bool) ServiceType {
	packagePath := t.PkgPath()
	if len(packagePath) == 0 {
		packagePath = "anonymous"
	}
	name := elemType.Name()
	if len(name) == 0 {
		name = t.String()
	}
	key := fmt.Sprintf("%s.%s", packagePath, name)
	if isList {
		key = fmt.Sprintf("%s.%s[]", packagePath, name)
	}
	serviceKey := NewServiceKey(key)
	return &serviceType{
		reflectedType: t,
		name:          name,
		packagePath:   packagePath,
		list:          isList,
		lookupKey:     serviceKey,
	}
}
