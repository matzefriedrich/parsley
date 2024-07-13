package tests

import (
	"context"
	"github.com/matzefriedrich/parsley/internal/core"
	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/resolving"
	"github.com/matzefriedrich/parsley/pkg/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Resolver_Resolve_returns_err_if_circular_dependency_detected(t *testing.T) {

	// Arrange
	registry := registration.NewServiceRegistry()
	_ = registration.RegisterTransient(registry, newFoo)
	_ = registration.RegisterTransient(registry, newBar)

	r := resolving.NewResolver(registry)

	scope := core.NewScopedContext(context.Background())

	// Act
	_, err := r.Resolve(scope, registration.ServiceType[Foo0]())

	// Assert
	assert.ErrorIs(t, err, types.ErrCircularDependencyDetected)
	assert.ErrorIs(t, err, types.ErrCannotBuildDependencyGraph)
}

type foo0 struct {
	bar Bar0
}

type Foo0 interface{}

type bar0 struct {
	foo Foo0
}

type Bar0 interface{}

func newFoo(bar Bar0) Foo0 {
	return &foo0{bar: bar}
}

func newBar(foo Foo0) Bar0 {
	return &bar0{foo: foo}
}
