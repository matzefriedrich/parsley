package registration

import (
	"testing"

	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/resolving"
	"github.com/matzefriedrich/parsley/pkg/types"
	"github.com/stretchr/testify/assert"
)

func Test_Register_struct_dependency(t *testing.T) {

	// Arrange
	registry := registration.NewServiceRegistry()

	config := someConfig{b: true}
	registryErr := registration.RegisterInstance[someConfig](registry, config)

	resolver := resolving.NewResolver(registry)
	ctx := resolving.NewScopedContext(t.Context())

	// Act
	actual, err := resolving.ResolveRequiredService[someConfig](ctx, resolver)

	// Arrange
	assert.NoError(t, registryErr)
	assert.NoError(t, err)
	assert.True(t, actual.b)
}

func Test_Register_service_with_struct_dependency(t *testing.T) {

	// Arrange
	registry := registration.NewServiceRegistry()
	registry.Register(newAppWithStructDependency, types.LifetimeTransient)

	config := someConfig{b: true}
	registryErr := registration.RegisterInstance[someConfig](registry, config)

	resolver := resolving.NewResolver(registry)
	ctx := resolving.NewScopedContext(t.Context())

	// Act
	actual, err := resolving.ResolveRequiredService[*appWithStructDependency](ctx, resolver)

	// Arrange
	assert.NoError(t, registryErr)
	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.True(t, actual.config.b)
}

func Test_Register_immutable_service_with_struct_dependency(t *testing.T) {

	// Arrange
	registry := registration.NewServiceRegistry()
	registry.Register(newImmutableAppWithStructDependency, types.LifetimeTransient)

	config := someConfig{b: true}
	registryErr := registration.RegisterInstance[someConfig](registry, config)

	resolver := resolving.NewResolver(registry)
	ctx := resolving.NewScopedContext(t.Context())

	// Act
	actual, err := resolving.ResolveRequiredService[appWithStructDependency](ctx, resolver)

	// Arrange
	assert.NoError(t, registryErr)
	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.True(t, actual.config.b)
}

type someConfig struct {
	b bool
}

type appWithStructDependency struct {
	config someConfig
}

func newAppWithStructDependency(config someConfig) *appWithStructDependency {
	return &appWithStructDependency{config: config}
}

func newImmutableAppWithStructDependency(config someConfig) appWithStructDependency {
	return appWithStructDependency{config: config}
}
