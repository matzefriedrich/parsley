package features

import (
	"fmt"
	"reflect"
)

type MethodInterceptor interface {
	Enter(target any, methodName string, parameters []ParameterInfo)
	Exit(target any, methodName string, returnValues []ReturnValueInfo)
	OnError(target any, methodName string, err error)
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

// AddResult Appends the result values to the current MethodCallContext instance.
func (c *MethodCallContext) AddResult(results ...any) {
	c.returnValues = append(c.returnValues, results...)
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

func (p *ProxyBase) InvokeMethodErrorInterceptors(callContext *MethodCallContext) {
	for _, next := range callContext.returnValues {
		err, ok := next.(error)
		if ok {
			for _, i := range p.interceptors {
				i.OnError(p.target, callContext.methodName, err)
			}
		}
	}
}

func (p *ProxyBase) InvokeEnterMethodInterceptors(callContext *MethodCallContext) {
	parameters := make([]ParameterInfo, 0, len(callContext.parameters))
	for name, next := range callContext.parameters {
		value := reflect.ValueOf(next)
		var parameterType reflect.Type
		switch value.Kind() {
		case reflect.Invalid:
			parameterType = nil
		default:
			parameterType = value.Type()
		}
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
		value := reflect.ValueOf(next)
		var valueType reflect.Type
		switch value.Kind() {
		case reflect.Invalid:
			valueType = nil
		default:
			valueType = value.Type()
		}
		returnValues = append(returnValues, ReturnValueInfo{
			value:     next,
			valueType: valueType,
		})
	}
	for _, i := range p.interceptors {
		i.Exit(p.target, callContext.methodName, returnValues)
	}
}

func NewProxyBase[T any](target T, interceptors []MethodInterceptor) ProxyBase {
	return ProxyBase{
		target:       target,
		interceptors: interceptors,
	}
}
