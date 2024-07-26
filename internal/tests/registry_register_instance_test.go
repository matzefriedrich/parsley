package tests

import (
	"context"
	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/resolving"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Registry_RegisterInstance_accepts_pointer(t *testing.T) {

	// Arrange
	registry := registration.NewServiceRegistry()

	options := NewOptions("value")

	// Act
	err := registration.RegisterInstance(registry, options)
	resolver := resolving.NewResolver(registry)
	actual, _ := resolving.ResolveRequiredService[*someOptions](resolver, resolving.NewScopedContext(context.Background()))

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, actual)
	assert.Equal(t, options.value, actual.value)
}

func Test_Registry_RegisterInstance_resolve_object_with_pointer_dependency(t *testing.T) {

	// Arrange
	registry := registration.NewServiceRegistry()

	options := NewOptions("value")
	_ = registration.RegisterInstance(registry, options)
	_ = registration.RegisterTransient(registry, NewOptionsConsumer)

	resolver := resolving.NewResolver(registry)

	// Act
	actual, _ := resolving.ResolveRequiredService[*optionsConsumer](resolver, resolving.NewScopedContext(context.Background()))

	// Assert
	assert.NotNil(t, actual)
}

type someOptions struct {
	value string
}

func NewOptions(value string) *someOptions {
	return &someOptions{value}
}

type optionsConsumer struct {
	options *someOptions
}

func NewOptionsConsumer(options *someOptions) *optionsConsumer {
	return &optionsConsumer{
		options: options,
	}
}
