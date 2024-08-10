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
    return &GreeterProxyImpl {
        ProxyBase: features.NewProxyBase(target, interceptors),
        target: target,
    }
}

func (__p *GreeterProxyImpl) SayHello(name string) (string, error){
    const methodName = "SayHello"
    parameters := map[string]interface{}{ 
		"name": name,
	}
	callContext := features.NewMethodCallContext(methodName, parameters)
	__p.InvokeEnterMethodInterceptors(callContext)
	defer func () {
	    __p.InvokeExitMethodInterceptors(callContext)
	    __p.InvokeMethodErrorInterceptors(callContext)
	}()
    result0, result1 := __p.target.SayHello(name)
    return result0, result1
}

var _ Greeter = &GreeterProxyImpl{}
