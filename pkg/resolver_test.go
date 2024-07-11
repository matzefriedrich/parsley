package pkg

import (
	"context"
	"github.com/matzefriedrich/parsley/internal"
	"github.com/matzefriedrich/parsley/pkg/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Resolver_Resolve_returns_err_if_circular_dependency_detected(t *testing.T) {

	// Arrange
	registry := NewServiceRegistry()
	_ = RegisterTransient(registry, newFoo)
	_ = RegisterTransient(registry, newBar)

	r := registry.BuildResolver()
	scope := internal.NewScopedContext(context.Background())

	// Act
	_, err := r.Resolve(scope, types.ServiceType[Foo0]())

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
