package features

import (
	"fmt"
	"reflect"
	"sort"
)

// Interceptor is a base interface type for defining interceptors that can be used to monitor or alter the behavior of other components.
type Interceptor interface {
	Name() string
	Position() int
}

// MethodInterceptor provides hooks to intercept method execution on a proxy object.
// It allows entering before method invocations, exiting after method executions, and handling errors during method execution for monitoring or altering behavior.
type MethodInterceptor interface {
	Interceptor
	Enter(target any, methodName string, parameters []ParameterInfo)
	Exit(target any, methodName string, returnValues []ReturnValueInfo)
	OnError(target any, methodName string, err error)
}

// InterceptorBase serves as a foundational structure for defining interceptors, managing essential data like name and position.
type InterceptorBase struct {
	name     string
	position int
}

// Name retrieves the name of the interceptor, which is useful for identification and debugging purposes.
func (i InterceptorBase) Name() string {
	return i.name
}

// Position returns the position of the interceptor, helping determine its order in processing flows within a system.
func (i InterceptorBase) Position() int {
	return i.position
}

// NewInterceptorBase creates a new instance of InterceptorBase with the specified name and position for managing interceptor metadata.
func NewInterceptorBase(name string, position int) InterceptorBase {
	return InterceptorBase{
		name:     name,
		position: position,
	}
}

// MethodCallContext captures the context of a method call, including method name, parameters, and return values.
type MethodCallContext struct {
	methodName   string
	parameters   map[string]interface{}
	returnValues []interface{}
}

// ParameterInfo represents information about a method parameter, including its value, type, and name.
// It is used in method interception where parameters need to be inspected or logged.
type ParameterInfo struct {
	value         interface{}
	parameterType reflect.Type
	name          string
}

// String returns a formatted string representation of the ParameterInfo, useful for logging and debugging purposes.
func (p ParameterInfo) String() string {
	return fmt.Sprintf("{%s (%s): %s}", p.name, p.parameterType, p.value)
}

// ReturnValueInfo represents the value and type information of a method's return value, used in method interception.
type ReturnValueInfo struct {
	value     interface{}
	valueType reflect.Type
}

// String returns a string representation of ReturnValueInfo, formatting the value and its type for debugging purposes.
func (r ReturnValueInfo) String() string {
	if r.value != nil {
		return fmt.Sprintf("{%s: %v}", r.valueType.String(), r.value)
	}
	return fmt.Sprintf("{%v}", r.value)
}

// NewMethodCallContext creates a new MethodCallContext instance with the provided method name and parameters.
func NewMethodCallContext(methodName string, parameters map[string]interface{}) *MethodCallContext {
	return &MethodCallContext{
		methodName:   methodName,
		parameters:   parameters,
		returnValues: make([]interface{}, 0),
	}
}

// ProxyBase facilitates method interception by allowing the inclusion of multiple interceptors to target method calls.
// Typically used to monitor, log, or modify behavior of an object's method execution.
type ProxyBase struct {
	target       any
	interceptors []MethodInterceptor
}

// InvokeMethodErrorInterceptors intercepts the return values of a method, checks for errors, and triggers OnError for registered interceptors.
func (p *ProxyBase) InvokeMethodErrorInterceptors(callContext *MethodCallContext, returnValues ...interface{}) {
	for _, next := range returnValues {
		callContext.returnValues = append(callContext.returnValues, next)
		err, ok := next.(error)
		if ok {
			wrapped := &proxyError{err: err}
			for _, i := range p.interceptors {
				i.OnError(p.target, callContext.methodName, wrapped)
			}
		}
	}
}

// InvokeEnterMethodInterceptors triggers the Enter method on all registered interceptors before the target method executes.
func (p *ProxyBase) InvokeEnterMethodInterceptors(callContext *MethodCallContext) {
	parameters := make([]ParameterInfo, 0, len(callContext.parameters))
	for name, next := range callContext.parameters {
		value, parameterType := reflectValueInfo(next)
		p := ParameterInfo{
			value:         value,
			parameterType: parameterType,
			name:          name,
		}
		parameters = append(parameters, p)
	}
	for _, i := range p.interceptors {
		i.Enter(p.target, callContext.methodName, parameters)
	}
}

// InvokeExitMethodInterceptors triggers the Exit method of all registered interceptors after the target method completes.
func (p *ProxyBase) InvokeExitMethodInterceptors(callContext *MethodCallContext) {
	returnValues := make([]ReturnValueInfo, 0, len(callContext.returnValues))
	for _, next := range callContext.returnValues {
		value, returnType := reflectValueInfo(next)
		returnValues = append(returnValues, ReturnValueInfo{
			value:     value,
			valueType: returnType,
		})
	}
	for _, i := range p.interceptors {
		i.Exit(p.target, callContext.methodName, returnValues)
	}
}

// NewProxyBase creates a ProxyBase instance with the provided target and a sorted list of method interceptors.
// Useful for setting up method interception on the target object.
func NewProxyBase[T any](target T, interceptors []MethodInterceptor) ProxyBase {
	sortedInterceptors := make([]MethodInterceptor, 0, len(interceptors))
	for _, interceptor := range interceptors {
		sortedInterceptors = append(sortedInterceptors, interceptor)
	}
	sort.Slice(sortedInterceptors, func(i, j int) bool {
		return sortedInterceptors[i].Position() < sortedInterceptors[j].Position()
	})
	return ProxyBase{
		target:       target,
		interceptors: sortedInterceptors,
	}
}

func reflectValueInfo(next interface{}) (reflect.Value, reflect.Type) {
	value := reflect.ValueOf(next)
	var parameterType reflect.Type
	switch value.Kind() {
	case reflect.Invalid:
		parameterType = nil
	default:
		parameterType = value.Type()
	}
	return value, parameterType
}

type proxyError struct {
	err error
}

// Error returns the error message encapsulated within the proxyError.
func (e proxyError) Error() string {
	return e.err.Error()
}

// Unwrap returns the underlying error encapsulated within proxyError for further inspection or handling.
func (e proxyError) Unwrap() error {
	return e.err
}

// Is checks if the provided error matches the encapsulated error by comparing their error messages.
func (e proxyError) Is(err error) bool {
	return e.err.Error() == err.Error()
}
