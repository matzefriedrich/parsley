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

func (m methodCallInterceptor) Enter(_ any, methodName string, parameters map[string]interface{}) {
	fmt.Println("Enter method: ", methodName)
}

func (m methodCallInterceptor) Exit(_ any, methodName string, returnValues []interface{}) {
	fmt.Println("Exit method: ", methodName)
}

func (m methodCallInterceptor) OnError(_ any, methodName string, err error) {
	fmt.Printf("OnError method: %s, Error: %v\n", methodName, err)
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
