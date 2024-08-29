package features

import (
	"fmt"
	"reflect"
	"sort"
)

type Interceptor interface {
	Name() string
	Position() int
}

type MethodInterceptor interface {
	Interceptor
	Enter(target any, methodName string, parameters []ParameterInfo)
	Exit(target any, methodName string, returnValues []ReturnValueInfo)
	OnError(target any, methodName string, err error)
}

type InterceptorBase struct {
	name     string
	position int
}

func (i InterceptorBase) Name() string {
	return i.name
}

func (i InterceptorBase) Position() int {
	return i.position
}

func NewInterceptorBase(name string, position int) InterceptorBase {
	return InterceptorBase{
		name:     name,
		position: position,
	}
}

type MethodCallContext struct {
	methodName   string
	parameters   map[string]interface{}
	returnValues []interface{}
}

type ParameterInfo struct {
	value         interface{}
	parameterType reflect.Type
	name          string
}

func (p ParameterInfo) String() string {
	return fmt.Sprintf("{%s (%s): %s}", p.name, p.parameterType, p.value)
}

type ReturnValueInfo struct {
	value     interface{}
	valueType reflect.Type
}

func (r ReturnValueInfo) String() string {
	if r.value != nil {
		return fmt.Sprintf("{%s: %v}", r.valueType.String(), r.value)
	}
	return fmt.Sprintf("{%v}", r.value)
}

func NewMethodCallContext(methodName string, parameters map[string]interface{}) *MethodCallContext {
	return &MethodCallContext{
		methodName:   methodName,
		parameters:   parameters,
		returnValues: make([]interface{}, 0),
	}
}

type ProxyBase struct {
	target       any
	interceptors []MethodInterceptor
}

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

func (e proxyError) Error() string {
	return e.err.Error()
}

func (e proxyError) Unwrap() error {
	return e.err
}

func (e proxyError) Is(err error) bool {
	return e.err.Error() == err.Error()
}
