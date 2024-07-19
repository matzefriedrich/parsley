package tests

import (
	"context"
	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/resolving"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Activate_resolve_unknown_service_type_using_resolve_options(t *testing.T) {

	// Arrange
	registry := registration.NewServiceRegistry()
	_ = registration.RegisterTransient(registry, newBar)

	sut := resolving.NewResolver(registry)
	scopedContext := resolving.NewScopedContext(context.Background())

	// Act
	actual, err := resolving.Activate[bar0](sut, scopedContext, newFooWithBar)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, actual)
}

type fooWithBar struct {
	bar bar0
}

type bar struct{}

type foo0 interface{}

type bar0 interface{}

func newFooWithBar(bar bar0) foo0 {
	return &fooWithBar{bar: bar}
}

func newBar() bar0 {
	return &bar{}
}
