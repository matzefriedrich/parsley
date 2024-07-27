package core

import (
	"errors"
	"fmt"
	"github.com/matzefriedrich/parsley/pkg/types"
	"reflect"
	"runtime"
	"strings"
)

const (
	ErrorNotAFunction                            = "not a function"
	ErrorReturnTypeHasToHaveExactlyOnReturnValue = "return type has to have exactly one return value"
)

var (
	ErrNotAFunction                            = errors.New(ErrorNotAFunction)
	ErrReturnTypeHasToHaveExactlyOnReturnValue = errors.New(ErrorReturnTypeHasToHaveExactlyOnReturnValue)
)

type functionInfo struct {
	reflectedFunctionValue reflect.Value
	funcType               reflect.Type
	returnType             types.ServiceType
	parameters             []types.FunctionParameterInfo
}

var _ types.FunctionInfo = &functionInfo{}

type functionParameterInfo struct {
	parameterType types.ServiceType
}

func (f functionParameterInfo) String() string {
	reflectedType := f.Type().ReflectedType()
	return fmt.Sprintf("%s", reflectedType.String())
}

func (f functionParameterInfo) Type() types.ServiceType {
	return f.parameterType
}

var _ types.FunctionParameterInfo = &functionParameterInfo{}

func ReflectFunctionInfoFrom(value reflect.Value) (types.FunctionInfo, error) {
	funcType := value.Type()
	if funcType.Kind() != reflect.Func {
		return nil, types.NewReflectionError(ErrorNotAFunction)
	}
	rt, err := returnType(funcType)
	if err != nil {
		return nil, err
	}
	parameters := parameterInfos(funcType)
	return &functionInfo{
		reflectedFunctionValue: value,
		funcType:               funcType,
		returnType:             rt,
		parameters:             parameters,
	}, nil
}

func (f functionInfo) Name() string {
	pointer := f.reflectedFunctionValue.Pointer()
	functionFromPointer := runtime.FuncForPC(pointer)
	if functionFromPointer != nil {
		return functionFromPointer.Name()
	}
	return ""
}

func (f functionInfo) Parameters() []types.FunctionParameterInfo {
	return f.parameters
}

func (f functionInfo) ParameterTypes() []types.ServiceType {
	parameterTypes := make([]types.ServiceType, len(f.parameters))
	for i, parameter := range f.parameters {
		parameterTypes[i] = parameter.Type()
	}
	return parameterTypes
}

func (f functionInfo) ReturnType() types.ServiceType {
	return f.returnType
}

func (f functionInfo) String() string {
	parameterTypeNames := make([]string, len(f.parameters))
	for _, t := range f.parameters {
		parameterTypeNames[0] = t.String()
	}
	reflectedReturnType := f.returnType.ReflectedType()
	funcTypeName := reflectedReturnType.String()
	return fmt.Sprintf("%s(%s) %s", f.Name(), strings.Join(parameterTypeNames, ","), funcTypeName)
}

func returnType(funcType reflect.Type) (types.ServiceType, error) {
	numReturnValues := funcType.NumOut()
	if numReturnValues != 1 {
		return nil, types.NewReflectionError(ErrorReturnTypeHasToHaveExactlyOnReturnValue)
	}
	serviceType := funcType.Out(0)
	return types.ServiceTypeFrom(serviceType), nil
}

func parameterInfos(funcType reflect.Type) []types.FunctionParameterInfo {
	parameters := make([]types.FunctionParameterInfo, 0)
	numParameters := funcType.NumIn()
	for i := 0; i < numParameters; i++ {
		parameterType := funcType.In(i)
		serviceType := types.ServiceTypeFrom(parameterType)
		p := functionParameterInfo{
			parameterType: serviceType,
		}
		parameters = append(parameters, p)
	}
	return parameters
}
