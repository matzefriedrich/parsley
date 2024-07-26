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

type someOptions struct {
	value string
}

func NewOptions(value string) *someOptions {
	return &someOptions{value}
}
