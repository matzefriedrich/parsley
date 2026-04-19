package features

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/matzefriedrich/parsley/pkg/features"
	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/resolving"
	"github.com/matzefriedrich/parsley/pkg/types"
)

func Test_Proxy_reflect_parameter_type_nil_value(t *testing.T) {

	// Arrange
	ctx := t.Context()
	collector := &nilParamCollector{}

	registry := registration.NewServiceRegistry()
	_ = registry.Register(newNilParamInterceptor(collector), types.LifetimeSingleton)
	_ = registry.Register(NewNilParamReproProxyImpl, types.LifetimeTransient)
	_ = registry.Register(newNilParamReproImpl, types.LifetimeTransient)
	_ = features.RegisterList[features.MethodInterceptor](ctx, registry)

	resolver := resolving.NewResolver(registry)
	resolverContext := resolving.NewScopedContext(ctx)

	// Act
	proxy, _ := resolving.ResolveRequiredService[NilParamReproProxy](resolverContext, resolver)
	proxy.SaySomething(nil)

	// Assert
	if collector.paramType == nil {
		t.Error("Parameter type was lost (nil) for nil parameter value")
	} else {
		fmt.Printf("Parameter type for nil error: %v\n", collector.paramType)
		if collector.paramType.String() != "error" {
			t.Errorf("Expected parameter type 'error', got '%v'", collector.paramType)
		}
	}
}

type nilParamCollector struct {
	paramType reflect.Type
}

type nilParamInterceptor struct {
	features.InterceptorBase
	collector *nilParamCollector
}

func (m nilParamInterceptor) Enter(_ any, methodName string, parameters []features.ParameterInfo) {
	if methodName == "SaySomething" {
		for _, p := range parameters {
			if p.Name() == "err" {
				m.collector.paramType = p.ParameterType()
			}
		}
	}
}

func (m nilParamInterceptor) Exit(_ any, _ string, _ []features.ReturnValueInfo) {}
func (m nilParamInterceptor) OnError(_ any, _ string, _ error)                   {}

var _ features.MethodInterceptor = &nilParamInterceptor{}

func newNilParamInterceptor(collector *nilParamCollector) func() features.MethodInterceptor {
	return func() features.MethodInterceptor {
		return &nilParamInterceptor{
			InterceptorBase: features.NewInterceptorBase("nil-param-interceptor", 0),
			collector:       collector,
		}
	}
}

type nilParamReproImpl struct{}

func (n *nilParamReproImpl) SaySomething(err error) {}

func newNilParamReproImpl() NilParamRepro {
	return &nilParamReproImpl{}
}
