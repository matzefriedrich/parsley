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

func (__p *GreeterProxyImpl) SayHello(name string) {
    const methodName = "SayHello"
    parameters := map[string]interface{}{ 
		"name": name,
	}
    __p.InvokeEnterMethodInterceptors(methodName, parameters)
    __p.target.SayHello(name)
    __p.InvokeExitMethodInterceptors(methodName, []interface{}{})
}

var _ Greeter = &GreeterProxyImpl{}
