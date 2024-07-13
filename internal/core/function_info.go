package core

import (
	"errors"
	"github.com/matzefriedrich/parsley/pkg/types"
	"reflect"
)

type functionInfo struct {
	funcType       reflect.Type
	returnType     reflect.Type
	parameterTypes []reflect.Type
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

func (f functionInfo) ParameterTypes() []reflect.Type {
	return f.parameterTypes
}

func (f functionInfo) ReturnType() reflect.Type {
	return f.returnType
}

func returnType(funcType reflect.Type) (reflect.Type, error) {
	numReturnValues := funcType.NumOut()
	if numReturnValues != 1 {
		return nil, errors.New("return type has to have exactly one return value")
	}
	serviceType := funcType.Out(0)
	return serviceType, nil
}

func parameterTypes(funcType reflect.Type) []reflect.Type {
	parameters := make([]reflect.Type, 0)
	numParameters := funcType.NumIn()
	for i := 0; i < numParameters; i++ {
		parameterType := funcType.In(i)
		parameters = append(parameters, parameterType)
	}
	return parameters
}
