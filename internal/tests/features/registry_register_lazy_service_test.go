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

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, actual)

	fooInstance0 := actual.Value(ctx)
	assert.NotNil(t, fooInstance0)

	fooInstance1 := actual.Value(ctx)
	assert.NotNil(t, fooInstance1)
}

func Test_Registry_register_lazy_service_with_dependency(t *testing.T) {

	// Arrange
	const expectedName = "Dependency"
	registry := registration.NewServiceRegistry()
	_ = registration.RegisterSingleton(registry, func() *dependency {
		return &dependency{Name: expectedName}
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
	assert.Equal(t, expectedName, serviceInstance.Dependency.Name)
}

func Test_Lazy_Value_returns_zero_value_on_activation_error(t *testing.T) {

	// Arrange
	registry := registration.NewServiceRegistry()

	// Register a lazy service that depends on something NOT registered
	activatorFunc := func(d *missingDependency) *serviceWithMissingDependency {
		return &serviceWithMissingDependency{Dependency: d}
	}
	_ = features.RegisterLazy[*serviceWithMissingDependency](registry, activatorFunc, types.LifetimeTransient)

	resolver := resolving.NewResolver(registry)
	ctx := resolving.NewScopedContext(t.Context())

	// Act
	lazyService, err := resolving.ResolveRequiredService[features.Lazy[*serviceWithMissingDependency]](ctx, resolver)
	assert.NoError(t, err)

	// This should fail to activate the underlying service because missingDependency is not registered
	serviceInstance := lazyService.Value(ctx)

	// Assert
	assert.Nil(t, serviceInstance)
}

func Test_Lazy_Value_is_thread_safe(t *testing.T) {

	// Arrange
	registry := registration.NewServiceRegistry()
	counter := 0

	activatorFunc := func() *foo {
		counter++
		return &foo{}
	}

	_ = features.RegisterLazy[*foo](registry, activatorFunc, types.LifetimeTransient)

	resolver := resolving.NewResolver(registry)
	ctx := resolving.NewScopedContext(t.Context())

	lazyService, _ := resolving.ResolveRequiredService[features.Lazy[*foo]](ctx, resolver)

	// Act
	const iterations = 100
	done := make(chan bool)
	for i := 0; i < iterations; i++ {
		go func() {
			_ = lazyService.Value(ctx)
			done <- true
		}()
	}

	for i := 0; i < iterations; i++ {
		<-done
	}

	// Assert
	assert.Equal(t, 1, counter)
}

type foo struct {
}

type dependency struct {
	Name string
}

type service struct {
	Dependency *dependency
}

type missingDependency struct {
}

type serviceWithMissingDependency struct {
	Dependency *missingDependency
}
