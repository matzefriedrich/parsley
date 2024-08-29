package features

import (
	"context"
	"errors"
	"fmt"
	"github.com/matzefriedrich/parsley/pkg/features"
	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/resolving"
	"github.com/matzefriedrich/parsley/pkg/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

type errorCollector struct {
	collected []error
}

type methodErrorInterceptor struct {
	features.InterceptorBase
	collector *errorCollector
}

var _ features.MethodInterceptor = &methodErrorInterceptor{}

func newMethodErrorInterceptor(collector *errorCollector) func() features.MethodInterceptor {
	return func() features.MethodInterceptor {
		return &methodErrorInterceptor{
			InterceptorBase: features.NewInterceptorBase("error-interceptor", 0),
			collector:       collector,
		}
	}
}

func (m methodErrorInterceptor) Enter(_ any, _ string, _ []features.ParameterInfo) {
}

func (m methodErrorInterceptor) Exit(_ any, _ string, _ []features.ReturnValueInfo) {
}

func (m methodErrorInterceptor) OnError(_ any, _ string, err error) {
	m.collector.collected = append(m.collector.collected, err)
}

func Test_Register_generated_proxy_type_handles_error(t *testing.T) {

	// Arrange
	collector := &errorCollector{collected: make([]error, 0)}

	registry := registration.NewServiceRegistry()
	registry.Register(newMethodErrorInterceptor(collector), types.LifetimeSingleton)
	registry.Register(NewGreeterProxyImpl, types.LifetimeTransient)
	registry.Register(newJohnGreeter, types.LifetimeTransient)
	features.RegisterList[features.MethodInterceptor](registry)

	resolver := resolving.NewResolver(registry)
	ctx := resolving.NewScopedContext(context.Background())

	// Act
	proxy, _ := resolving.ResolveRequiredService[GreeterProxy](resolver, ctx)
	msg, _ := proxy.SayHello("Jane")
	fmt.Println(msg)

	// Assert
	err := collector.collected[0]
	assert.Error(t, err)
	assert.ErrorIs(t, err, errors.New("name is not John"))
}

type johnGreeter struct {
}

func (g johnGreeter) SayNothing() {
}

func (g johnGreeter) SayHello(name string) (string, error) {
	if name != "John" {
		return "", fmt.Errorf("name is not John")
	}
	return "Hello " + name, nil
}

var _ Greeter = &johnGreeter{}

func newJohnGreeter() Greeter {
	return &johnGreeter{}
}
