package features

import (
    "github.com/matzefriedrich/parsley/pkg/features"
)

type GreeterProxyImpl struct {
    features.ProxyBase
    target Greeter
}

type GreeterProxy interface {
    Greeter
}

func NewGreeterProxyImpl(target Greeter, interceptors []features.MethodInterceptor) GreeterProxy {
    return &GreeterProxyImpl{
        ProxyBase: features.NewProxyBase(target, interceptors),
        target:    target,
    }
}

func (p *GreeterProxyImpl) SayHello(name string) (string, error) {

    const methodName = "SayHello"
    parameters := map[string]interface{}{ 
		"name": name,
	}

	callContext := features.NewMethodCallContext(methodName, parameters)
	p.InvokeEnterMethodInterceptors(callContext)
	defer func() {
	    p.InvokeExitMethodInterceptors(callContext)
	}()
    
    result0, result1 := p.target.SayHello(name)
    p.InvokeMethodErrorInterceptors(callContext, result0, result1)
    return result0, result1
}

func (p *GreeterProxyImpl) SayNothing()  {

    const methodName = "SayNothing"
    parameters := map[string]interface{}{ 	}

	callContext := features.NewMethodCallContext(methodName, parameters)
	p.InvokeEnterMethodInterceptors(callContext)
	defer func() {
	    p.InvokeExitMethodInterceptors(callContext)
	}()
    
    p.target.SayNothing()
}

var _ Greeter = &GreeterProxyImpl{}
