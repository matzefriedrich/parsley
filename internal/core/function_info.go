package core

import (
	"errors"
	"fmt"
	"github.com/matzefriedrich/parsley/pkg/types"
	"reflect"
	"strings"
)

type functionInfo struct {
	funcType       reflect.Type
	returnType     types.ServiceType
	parameterTypes []types.ServiceType
}

var _ types.FunctionInfo = &functionInfo{}

func ReflectFunctionInfoFrom(value reflect.Value) (types.FunctionInfo, error) {
	funcType := value.Type()
	if funcType.Kind() != reflect.Func {
		return nil, errors.New("not a function")
	}
	rt, err := returnType(funcType)
	if err != nil {
		return nil, err
	}
	parameters := parameterTypes(funcType)
	return &functionInfo{
		funcType:       funcType,
		returnType:     rt,
		parameterTypes: parameters,
	}, nil
}

func (f functionInfo) Name() string {
	return f.funcType.Name()
}

func (f functionInfo) ParameterTypes() []types.ServiceType {
	return f.parameterTypes
}

func (f functionInfo) ReturnType() types.ServiceType {
	return f.returnType
}

func (f functionInfo) String() string {
	parameterTypeNames := make([]string, len(f.parameterTypes))
	for _, t := range f.parameterTypes {
		parameterTypeNames = append(parameterTypeNames, t.Name())
	}
	funcTypeName := f.returnType.ReflectedType().Elem().String()
	return fmt.Sprintf("f(%s) %s", strings.Join(parameterTypeNames, ","), funcTypeName)
}

func returnType(funcType reflect.Type) (types.ServiceType, error) {
	numReturnValues := funcType.NumOut()
	if numReturnValues != 1 {
		return nil, errors.New("return type has to have exactly one return value")
	}
	serviceType := funcType.Out(0)
	return types.ServiceTypeFrom(serviceType), nil
}

func parameterTypes(funcType reflect.Type) []types.ServiceType {
	parameters := make([]types.ServiceType, 0)
	numParameters := funcType.NumIn()
	for i := 0; i < numParameters; i++ {
		parameterType := funcType.In(i)
		serviceType := types.ServiceTypeFrom(parameterType)
		parameters = append(parameters, serviceType)
	}
	return parameters
}
