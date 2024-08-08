package features

import (
	"context"
	"fmt"
	"github.com/matzefriedrich/parsley/pkg/features"
	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/resolving"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ProxyGenerator_generate_greeter_proxy(t *testing.T) {

	// Arrange
	registry := registration.NewServiceRegistry()
	registration.RegisterTransient(registry, newGreeter)
	registration.RegisterTransient(registry, NewGreeterProxyImpl)
	registration.RegisterInstance[features.MethodInterceptor](registry, &testMethodInterceptor{})

	resolver := resolving.NewResolver(registry)
	ctx := resolving.NewScopedContext(context.Background())

	// Act
	proxy, _ := resolving.ResolveRequiredService[GreeterProxy](resolver, ctx)

	// Assert
	assert.NotNil(t, proxy)

	proxy.SayHello("John")
}

type testMethodInterceptor struct {
}

func (t *testMethodInterceptor) Enter(target any, methodName string, parameters map[string]interface{}) {
	fmt.Printf("Enter method: %s\n", methodName)
}

func (t *testMethodInterceptor) Exit(target any, methodName string, returnValues []interface{}) {
	fmt.Printf("Exit method: %s\n", methodName)
}

func (t *testMethodInterceptor) OnError(target any, methodName string, err error) {
}

type greeter struct {
}

func (g *greeter) SayHello(name string) {
	fmt.Println("Hello " + name)
}

func newGreeter() Greeter {
	return &greeter{}
}
