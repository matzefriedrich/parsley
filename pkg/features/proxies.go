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
	methodName     string
	parameterNames []string
	parameters     map[string]interface{}
	returnNames    []string
	returnValues   []ReturnValueInfo
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
	return fmt.Sprintf("{%s (%s): %v}", p.name, p.parameterType, p.value)
}

// Name returns the parameter name.
func (p ParameterInfo) Name() string {
	return p.name
}

// Value returns the value of the parameter.
func (p ParameterInfo) Value() interface{} {
	return p.value
}

// ParameterType retrieves the reflected type of the parameter.
func (p ParameterInfo) ParameterType() reflect.Type {
	return p.parameterType
}

// ReturnValueInfo represents the value and type information of a method's return value, used in method interception.
type ReturnValueInfo struct {
	name      string
	value     interface{}
	valueType reflect.Type
}

// Name returns the return value name.
func (r ReturnValueInfo) Name() string {
	return r.name
}

// Value returns the value stored in the ReturnValueInfo instance.
func (r ReturnValueInfo) Value() interface{} {
	return r.value
}

// ValueType retrieves the return value type.
func (r ReturnValueInfo) ValueType() reflect.Type {
	return r.valueType
}

// String returns a string representation of ReturnValueInfo, formatting the value and its type for debugging purposes.
func (r ReturnValueInfo) String() string {
	valueTypeName := r.valueType.String()
	if r.value != nil {
		return fmt.Sprintf("{%s (%s): %v}", r.name, valueTypeName, r.value)
	}
	return fmt.Sprintf("{%s (%s)}: nil", r.name, valueTypeName)
}

// NewReturnValueInfo creates a new ReturnValueInfo object.
func NewReturnValueInfo(name string, value any, valueType reflect.Type) ReturnValueInfo {
	return ReturnValueInfo{
		name:      name,
		value:     value,
		valueType: valueType,
	}
}

// NewMethodCallContext creates a new MethodCallContext instance with the provided method name, parameters, and return value names.
func NewMethodCallContext(methodName string, parameterNames []string, parameters map[string]interface{}, returnNames ...string) *MethodCallContext {
	return &MethodCallContext{
		methodName:     methodName,
		parameterNames: parameterNames,
		parameters:     parameters,
		returnNames:    returnNames,
		returnValues:   make([]ReturnValueInfo, 0),
	}
}

// ProxyBase facilitates method interception by allowing the inclusion of multiple interceptors to target method calls.
// Typically used to monitor, log, or modify the behavior of an object's method execution.
type ProxyBase struct {
	target       any
	targetType   reflect.Type
	interceptors []MethodInterceptor
}

// InvokeMethodErrorInterceptors intercepts the return values of a method, checks for errors, and triggers OnError for registered interceptors.
func (p *ProxyBase) InvokeMethodErrorInterceptors(callContext *MethodCallContext, returnValues ...any) {
	method, ok := p.targetType.MethodByName(callContext.methodName)
	if !ok {
		return
	}

	for i, next := range returnValues {
		valueType := method.Type.Out(i)
		name := fmt.Sprintf("result%d", i)
		if i < len(callContext.returnNames) {
			name = callContext.returnNames[i]
		}
		info := NewReturnValueInfo(name, next, valueType)
		callContext.returnValues = append(callContext.returnValues, info)
		if next == nil {
			continue
		}
		err, ok := next.(error)
		if ok {
			wrapped := &proxyError{err: err}
			for _, interceptor := range p.interceptors {
				interceptor.OnError(p.target, callContext.methodName, wrapped)
			}
		}
	}
}

// InvokeEnterMethodInterceptors triggers the Enter method on all registered interceptors before the target method executes.
func (p *ProxyBase) InvokeEnterMethodInterceptors(callContext *MethodCallContext) {
	method, ok := p.targetType.MethodByName(callContext.methodName)
	if !ok {
		return
	}

	parameters := make([]ParameterInfo, 0, len(callContext.parameterNames))
	for i, name := range callContext.parameterNames {
		value := callContext.parameters[name]
		parameterType := method.Type.In(i)
		pInfo := ParameterInfo{
			value:         value,
			parameterType: parameterType,
			name:          name,
		}
		parameters = append(parameters, pInfo)
	}
	for _, i := range p.interceptors {
		i.Enter(p.target, callContext.methodName, parameters)
	}
}

// InvokeExitMethodInterceptors triggers the Exit method of all registered interceptors after the target method completes.
func (p *ProxyBase) InvokeExitMethodInterceptors(callContext *MethodCallContext) {
	for _, i := range p.interceptors {
		i.Exit(p.target, callContext.methodName, callContext.returnValues)
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
		targetType:   reflect.TypeOf((*T)(nil)).Elem(),
		interceptors: sortedInterceptors,
	}
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
