package features

import (
	"testing"

	"github.com/matzefriedrich/parsley/pkg/features"
	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/resolving"
	"github.com/matzefriedrich/parsley/pkg/types"
	"github.com/stretchr/testify/assert"
)

func Test_Registry_register_lazy_service_type_factory(t *testing.T) {

	// Arrange
	registry := registration.NewServiceRegistry()
	_ = features.RegisterLazy[*foo](registry, func() *foo {
		return &foo{}
	}, types.LifetimeTransient)

	resolver := resolving.NewResolver(registry)
	ctx := resolving.NewScopedContext(t.Context())

	// Act
	actual, err := resolving.ResolveRequiredService[features.Lazy[*foo]](ctx, resolver)
	fooInstance0 := actual.Value(ctx)
	fooInstance1 := actual.Value(ctx)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.NotNil(t, fooInstance0)
	assert.NotNil(t, fooInstance1)
}

func Test_Registry_register_lazy_service_with_dependency(t *testing.T) {

	// Arrange
	registry := registration.NewServiceRegistry()
	_ = registration.RegisterSingleton(registry, func() *dependency {
		return &dependency{Name: "Dependency"}
	})

	_ = features.RegisterLazy[*service](registry, func(d *dependency) *service {
		return &service{Dependency: d}
	}, types.LifetimeTransient)

	resolver := resolving.NewResolver(registry)
	ctx := resolving.NewScopedContext(t.Context())

	// Act
	lazyService, err := resolving.ResolveRequiredService[features.Lazy[*service]](ctx, resolver)
	assert.NoError(t, err)

	serviceInstance := lazyService.Value(ctx)

	// Assert
	assert.NotNil(t, serviceInstance)
	assert.NotNil(t, serviceInstance.Dependency)
	assert.Equal(t, "Dependency", serviceInstance.Dependency.Name)
}

type foo struct {
}

type dependency struct {
	Name string
}

type service struct {
	Dependency *dependency
}
