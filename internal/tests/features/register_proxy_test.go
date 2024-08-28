package features

import (
	"context"
	"fmt"
	"github.com/matzefriedrich/parsley/pkg/features"
	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/resolving"
	"github.com/matzefriedrich/parsley/pkg/types"
	"testing"
)

func Test_Register_generated_proxy_type(t *testing.T) {

	// Arrange
	registry := registration.NewServiceRegistry()
	registry.Register(newMethodCallInterceptor, types.LifetimeSingleton)
	registry.Register(NewGreeterProxyImpl, types.LifetimeTransient)
	registry.Register(newGreeter, types.LifetimeTransient)
	features.RegisterList[features.MethodInterceptor](registry)

	resolver := resolving.NewResolver(registry)
	ctx := resolving.NewScopedContext(context.Background())

	// Act
	proxy, _ := resolving.ResolveRequiredService[GreeterProxy](resolver, ctx)
	msg, _ := proxy.SayHello("John")
	fmt.Println(msg)

	// Assert

}

type methodCallInterceptor struct {
}

func (m methodCallInterceptor) Enter(_ any, methodName string, parameters []features.ParameterInfo) {
	fmt.Println("Enter method: ", methodName)
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

func newMethodCallInterceptor() features.MethodInterceptor {
	return &methodCallInterceptor{}
}

type greeter struct {
}

func (g greeter) SayHello(name string) (string, error) {
	return "Hello " + name, nil
}

var _ Greeter = &greeter{}

func newGreeter() Greeter {
	return &greeter{}
}
