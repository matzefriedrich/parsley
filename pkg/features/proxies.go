package features

type MethodInterceptor interface {
	Enter(target any, methodName string, parameters map[string]interface{})
	Exit(target any, methodName string, returnValues []interface{})
	OnError(target any, methodName string, err error)
}

type MethodCallContext struct {
	methodName   string
	parameters   map[string]interface{}
	returnValues []interface{}
	err          error
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
	if callContext.err == nil {
		return
	}
	for _, i := range p.interceptors {
		i.OnError(p.target, callContext.methodName, callContext.err)
	}
}

func (p *ProxyBase) InvokeEnterMethodInterceptors(callContext *MethodCallContext) {
	for _, i := range p.interceptors {
		i.Enter(p.target, callContext.methodName, callContext.parameters)
	}
}

func (p *ProxyBase) InvokeExitMethodInterceptors(callContext *MethodCallContext) {
	for _, i := range p.interceptors {
		i.Exit(p.target, callContext.methodName, callContext.returnValues)
	}
}

func NewProxyBase[T any](target T, interceptors []MethodInterceptor) ProxyBase {
	return ProxyBase{
		target:       target,
		interceptors: interceptors,
	}
}
