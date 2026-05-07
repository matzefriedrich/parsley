package core

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/matzefriedrich/parsley/pkg/types"
)

const (
	ErrorNotAFunction                            = "not a function"
	ErrorReturnTypeHasToHaveExactlyOnReturnValue = "return type has to have exactly one return value"
	ErrorSecondReturnTypeIsNotErr                = "second return type is not an error type"
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
	hasErrorReturn         bool
	hasContextParameter    bool
}

var _ types.FunctionInfo = &functionInfo{}

type functionParameterInfo struct {
	parameterType types.ServiceType
}

// String returns the string representation of the reflected type of the function parameter.
func (f functionParameterInfo) String() string {
	reflectedType := f.Type().ReflectedType()
	return reflectedType.String()
}

// Type returns the ServiceType of the function parameter, which provides meta information like name, package path, and reflect type.
func (f functionParameterInfo) Type() types.ServiceType {
	return f.parameterType
}

var _ types.FunctionParameterInfo = &functionParameterInfo{}

// ReflectFunctionInfoFrom retrieves metadata about a given function using reflection, providing details about the function's type,
// return type, parameter types, and other characteristics. Ensures that the value provided is a valid function and returns
// appropriate errors if it is not. Useful for dynamically inspecting and working with functions.
func ReflectFunctionInfoFrom(value reflect.Value) (types.FunctionInfo, error) {
	funcType := value.Type()
	if funcType.Kind() != reflect.Func {
		return nil, types.NewReflectionError(ErrorNotAFunction)
	}
	rt, hasError, err := returnType(funcType)
	if err != nil {
		return nil, err
	}
	parameters := parameterInfos(funcType)
	hasContext := false
	if len(parameters) > 0 {
		firstParamType := parameters[0].Type().ReflectedType()
		contextType := reflect.TypeOf((*context.Context)(nil)).Elem()
		if firstParamType.Implements(contextType) {
			hasContext = true
		}
	}
	return &functionInfo{
		reflectedFunctionValue: value,
		funcType:               funcType,
		returnType:             rt,
		parameters:             parameters,
		hasErrorReturn:         hasError,
		hasContextParameter:    hasContext,
	}, nil
}

// Name retrieves the name of the reflected function.
func (f functionInfo) Name() string {
	pointer := f.reflectedFunctionValue.Pointer()
	functionFromPointer := runtime.FuncForPC(pointer)
	if functionFromPointer != nil {
		return functionFromPointer.Name()
	}
	return ""
}

// Parameters Returns a slice of FunctionParameterInfo representing the parameters of the function.
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

func (f functionInfo) HasErrorReturn() bool {
	return f.hasErrorReturn
}

func (f functionInfo) ExpectsContextParameter() bool {
	return f.hasContextParameter
}

// ReturnType returns the type of the service returned by the function.
func (f functionInfo) ReturnType() types.ServiceType {
	return f.returnType
}

// String returns a string representation of the function's signature, including its name, parameters, and return type.
func (f functionInfo) String() string {
	parameterTypeNames := make([]string, len(f.parameters))
	for i, t := range f.parameters {
		parameterTypeNames[i] = t.String()
	}
	reflectedReturnType := f.returnType.ReflectedType()
	funcTypeName := reflectedReturnType.String()
	return fmt.Sprintf("%s(%s) %s", f.Name(), strings.Join(parameterTypeNames, ","), funcTypeName)
}

func returnType(funcType reflect.Type) (types.ServiceType, bool, error) {
	numReturnValues := funcType.NumOut()
	const (
		serviceTypeIndex = 0
		errorTypeIndex   = 1
	)
	if numReturnValues == 1 {
		serviceType := funcType.Out(serviceTypeIndex)
		return types.ServiceTypeFrom(serviceType), false, nil
	}
	if numReturnValues == 2 {
		serviceType := funcType.Out(serviceTypeIndex)
		errorType := funcType.Out(errorTypeIndex)
		if isErrorType(errorType) {
			return types.ServiceTypeFrom(serviceType), true, nil
		}
		return nil, false, types.NewReflectionError(ErrorSecondReturnTypeIsNotErr)
	}
	return nil, false, types.NewReflectionError(ErrorReturnTypeHasToHaveExactlyOnReturnValue)
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

func isErrorType(t reflect.Type) bool {
	errorType := reflect.TypeOf((*error)(nil))
	return t.Implements(errorType.Elem())
}
