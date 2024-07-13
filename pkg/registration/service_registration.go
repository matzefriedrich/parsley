package registration

import (
	"errors"
	"fmt"
	"github.com/matzefriedrich/parsley/internal/core"
	"github.com/matzefriedrich/parsley/pkg/types"
	"reflect"
	"strings"
)

type serviceRegistration struct {
	id            uint64
	serviceType   typeInfo
	activatorFunc reflect.Value
	parameters    []typeInfo
	lifetimeScope types.LifetimeScope
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

func (s *serviceRegistration) InvokeActivator(params ...interface{}) (interface{}, error) {
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

func (s *serviceRegistration) Id() uint64 {
	return s.id
}

func (s *serviceRegistration) SetId(id uint64) error {
	if s.id != 0 {
		return errors.New("the id cannot be changed once set")
	}
	s.id = id
	return nil
}

func (s *serviceRegistration) LifetimeScope() types.LifetimeScope {
	return s.lifetimeScope
}

func (s *serviceRegistration) RequiredServiceTypes() []reflect.Type {
	requiredTypes := make([]reflect.Type, len(s.parameters))
	for i, p := range s.parameters {
		requiredTypes[i] = p.t
	}
	return requiredTypes
}

func (s *serviceRegistration) ServiceType() reflect.Type {
	return s.serviceType.t
}

func (s *serviceRegistration) String() string {

	buffer := strings.Builder{}
	parameterTypesNames := make([]string, 0)
	for _, parameterType := range s.parameters {
		parameterTypesNames = append(parameterTypesNames, parameterType.name)
	}

	buffer.WriteString(fmt.Sprintf("%s(%s)", s.serviceType.name, strings.Join(parameterTypesNames, ", ")))
	return buffer.String()
}

func CreateServiceRegistration(activatorFunc any, lifetimeScope types.LifetimeScope) (types.ServiceRegistrationSetup, error) {
	value := reflect.ValueOf(activatorFunc)

	info, err := core.ReflectFunctionInfoFrom(value)
	if err != nil {
		return nil, types.NewRegistryError(types.ErrorRequiresFunctionValue, types.WithCause(err))
	}

	serviceType := info.ReturnType()
	switch serviceType.Kind() {
	case reflect.Func:
		return newServiceRegistration(serviceType, lifetimeScope, value), nil
	case reflect.Interface:
		requiredTypes := info.ParameterTypes()
		return newServiceRegistration(serviceType, lifetimeScope, value, requiredTypes...), nil
	default:
		return nil, types.NewRegistryError(types.ErrorActivatorFunctionInvalidReturnType)
	}
}

func newServiceRegistration(serviceType reflect.Type, scope types.LifetimeScope, activatorFunc reflect.Value, parameters ...reflect.Type) *serviceRegistration {
	parameterTypeInfos := make([]typeInfo, len(parameters))
	for i, p := range parameters {
		parameterTypeInfos[i] = newTypeInfo(p)
	}
	return &serviceRegistration{
		serviceType:   newTypeInfo(serviceType),
		activatorFunc: activatorFunc,
		parameters:    parameterTypeInfos,
		lifetimeScope: scope,
	}
}

var _ types.ServiceRegistration = &serviceRegistration{}
var _ types.ServiceRegistrationSetup = &serviceRegistration{}
