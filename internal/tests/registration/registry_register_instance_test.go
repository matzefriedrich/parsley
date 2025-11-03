package registration

import (
	"testing"

	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/resolving"
	"github.com/stretchr/testify/assert"
)

func Test_Registry_RegisterInstance_accepts_pointer(t *testing.T) {

	// Arrange
	registry := registration.NewServiceRegistry()

	options := newOptions("value")

	// Act
	err := registration.RegisterInstance(registry, options)
	resolver := resolving.NewResolver(registry)
	resolverContext := resolving.NewScopedContext(t.Context())
	actual, _ := resolving.ResolveRequiredService[*someOptions](resolverContext, resolver)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, options.value, actual.value)
}

func Test_Registry_RegisterInstance_resolve_object_with_pointer_dependency(t *testing.T) {

	// Arrange
	registry := registration.NewServiceRegistry()

	options := newOptions("value")
	_ = registration.RegisterInstance(registry, options)
	_ = registration.RegisterTransient(registry, newOptionsConsumer)

	resolver := resolving.NewResolver(registry)
	resolverContext := resolving.NewScopedContext(t.Context())

	// Act
	actual, _ := resolving.ResolveRequiredService[*optionsConsumer](resolverContext, resolver)

	// Assert
	assert.NotNil(t, actual)
}

type someOptions struct {
	value string
}

func newOptions(value string) *someOptions {
	return &someOptions{value}
}

type optionsConsumer struct {
	options *someOptions
}

func newOptionsConsumer(options *someOptions) *optionsConsumer {
	return &optionsConsumer{
		options: options,
	}
}
