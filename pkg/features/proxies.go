package features

type MethodInterceptor interface {
	Enter(target any, methodName string, parameters map[string]interface{})
	Exit(target any, methodName string, returnValues []interface{})
	OnError(target any, methodName string, err error)
}

type ProxyBase struct {
	target       any
	interceptors []MethodInterceptor
}

func (p *ProxyBase) InvokeMethodErrorInterceptors(methodName string, err error) {
	for _, i := range p.interceptors {
		i.OnError(p.target, methodName, err)
	}
}

func (p *ProxyBase) InvokeEnterMethodInterceptors(methodName string, parameters map[string]interface{}) {
	for _, i := range p.interceptors {
		i.Enter(p.target, methodName, parameters)
	}
}

func (p *ProxyBase) InvokeExitMethodInterceptors(methodName string, returnValues []interface{}) {
	for _, i := range p.interceptors {
		i.Exit(p.target, methodName, returnValues)
	}
}

func NewProxyBase[T any](target T, interceptors []MethodInterceptor) ProxyBase {
	return ProxyBase{
		target:       target,
		interceptors: interceptors,
	}
}
