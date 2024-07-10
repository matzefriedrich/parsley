package pkg

import (
	"fmt"
	"github.com/matzefriedrich/parsley/pkg/types"
	"reflect"
	"strings"
)

type serviceRegistration struct {
	id            uint64
	serviceType   typeInfo
	activatorFunc reflect.Value
	parameters    []typeInfo
}

type typeInfo struct {
	t    reflect.Type
	name string
}

func newTypeInfo(t reflect.Type) typeInfo {
	return typeInfo{
		t:    t,
		name: t.Name(),
	}
}

func (s serviceRegistration) InvokeActivator(params ...interface{}) (interface{}, error) {
	var values []reflect.Value
	if len(params) > 0 {
		values = make([]reflect.Value, len(params))
		for i, p := range params {
			values[i] = reflect.ValueOf(p)
		}
	}
	result := s.activatorFunc.Call(values)
	if len(result) != 1 {
		return nil, fmt.Errorf("activator function returned %d values", len(result))
	}
	serviceInstance := result[0]
	return serviceInstance.Interface(), nil
}

func (s serviceRegistration) Id() uint64 {
	return s.id
}

func (s serviceRegistration) ServiceType() reflect.Type {
	return s.serviceType.t
}

func (s serviceRegistration) RequiredServiceTypes() []reflect.Type {
	types := make([]reflect.Type, len(s.parameters))
	for i, p := range s.parameters {
		types[i] = p.t
	}
	return types
}

func newServiceRegistration(serviceType reflect.Type, activatorFunc reflect.Value, parameters ...reflect.Type) *serviceRegistration {
	parameterTypeInfos := make([]typeInfo, len(parameters))
	for i, p := range parameters {
		parameterTypeInfos[i] = newTypeInfo(p)
	}
	return &serviceRegistration{
		serviceType:   newTypeInfo(serviceType),
		activatorFunc: activatorFunc,
		parameters:    parameterTypeInfos,
	}
}

func (s serviceRegistration) String() string {

	buffer := strings.Builder{}
	parameterTypesNames := make([]string, 0)
	for _, parameterType := range s.parameters {
		parameterTypesNames = append(parameterTypesNames, parameterType.name)
	}

	buffer.WriteString(fmt.Sprintf("%s(%s)", s.serviceType.name, strings.Join(parameterTypesNames, ", ")))
	return buffer.String()
}

var _ types.ServiceRegistration = &serviceRegistration{}
