package features

import (
	"fmt"
	"testing"

	"github.com/matzefriedrich/parsley/pkg/features"
	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/resolving"
	"github.com/matzefriedrich/parsley/pkg/types"
)

func Test_Register_generated_proxy_type(t *testing.T) {

	// Arrange
	ctx := t.Context()
	collector := &callCollector{methods: make([]string, 0)}

	registry := registration.NewServiceRegistry()
	registry.Register(newMethodCallInterceptor(collector), types.LifetimeSingleton)
	registry.Register(NewGreeterProxyImpl, types.LifetimeTransient)
	registry.Register(newGreeter, types.LifetimeTransient)
	features.RegisterList[features.MethodInterceptor](ctx, registry)

	resolver := resolving.NewResolver(registry)
	resolverContext := resolving.NewScopedContext(ctx)

	// Act
	proxy, _ := resolving.ResolveRequiredService[GreeterProxy](resolverContext, resolver)
	msg, _ := proxy.SayHello("John", false)
	fmt.Println(msg)

	// Assert

}

type callCollector struct {
	methods []string
}

func (c *callCollector) Verify(method string) bool {
	for _, m := range c.methods {
		if m == method {
			return true
		}
	}
	return false
}

type methodCallInterceptor struct {
	features.InterceptorBase
	calls *callCollector
}

func (m methodCallInterceptor) Enter(_ any, methodName string, parameters []features.ParameterInfo) {
	switch methodName {
	case "SayHello":
	}
	fmt.Println("Enter method: ", methodName)
	m.calls.methods = append(m.calls.methods, methodName)
	for _, parameterValue := range parameters {
		fmt.Printf("\t%s\n", parameterValue)
	}
}

func (m methodCallInterceptor) Exit(_ any, methodName string, returnValues []features.ReturnValueInfo) {
	fmt.Println("Exit method: ", methodName)
	for _, value := range returnValues {
		fmt.Printf("\tResult: %s\n", value)
	}
}

func (m methodCallInterceptor) OnError(_ any, _ string, _ error) {
}

var _ features.MethodInterceptor = &methodCallInterceptor{}

func newMethodCallInterceptor(collector *callCollector) func() features.MethodInterceptor {
	return func() features.MethodInterceptor {
		return &methodCallInterceptor{
			InterceptorBase: features.NewInterceptorBase("call-interceptor", 0),
			calls:           collector,
		}
	}
}

type greeter struct {
}

func (g greeter) SayNothing() {
}

func (g greeter) SayHello(name string, _ bool) (string, error) {
	return "Hello " + name, nil
}

var _ Greeter = &greeter{}

func newGreeter() Greeter {
	return &greeter{}
}
