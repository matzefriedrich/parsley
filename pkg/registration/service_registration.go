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
	t    types.ServiceType
	name string
}

func newTypeInfo(t types.ServiceType) typeInfo {
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

func (s *serviceRegistration) IsSame(other types.ServiceRegistration) bool {
	sr, ok := other.(*serviceRegistration)
	return ok && s.activatorFunc.Pointer() == sr.activatorFunc.Pointer()
}

func (s *serviceRegistration) LifetimeScope() types.LifetimeScope {
	return s.lifetimeScope
}

func (s *serviceRegistration) RequiredServiceTypes() []types.ServiceType {
	requiredTypes := make([]types.ServiceType, len(s.parameters))
	for i, p := range s.parameters {
		requiredTypes[i] = p.t
	}
	return requiredTypes
}

func (s *serviceRegistration) ServiceType() types.ServiceType {
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
	switch serviceType.ReflectedType().Kind() {
	case reflect.Func:
		fallthrough
	case reflect.Pointer:
		fallthrough
	case reflect.Interface:
		requiredTypes := info.ParameterTypes()
		return newServiceRegistration(serviceType, lifetimeScope, value, requiredTypes...), nil
	default:
		return nil, types.NewRegistryError(types.ErrorActivatorFunctionInvalidReturnType)
	}
}

func newServiceRegistration(serviceType types.ServiceType, scope types.LifetimeScope, activatorFunc reflect.Value, parameters ...types.ServiceType) *serviceRegistration {
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
