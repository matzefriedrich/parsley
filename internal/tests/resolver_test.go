package tests

import (
	"context"
	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/resolving"
	"github.com/matzefriedrich/parsley/pkg/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Resolver_Resolve_returns_err_if_circular_dependency_detected(t *testing.T) {

	// Arrange
	registry := registration.NewServiceRegistry()
	_ = registration.RegisterTransient(registry, newFooRequiringBar)
	_ = registration.RegisterTransient(registry, newBarRequiringFoo)

	r := resolving.NewResolver(registry)

	scope := resolving.NewScopedContext(context.Background())

	// Act
	_, err := r.Resolve(scope, registration.ServiceType[fooBar]())

	// Assert
	assert.ErrorIs(t, err, types.ErrCircularDependencyDetected)
	assert.ErrorIs(t, err, types.ErrCannotBuildDependencyGraph)
}

type fooRequiresBar struct {
	bar barFoo
}

type fooBar interface{}

type barRequiresFoo struct {
	foo fooBar
}

type barFoo interface{}

func newFooRequiringBar(bar barFoo) fooBar {
	return &fooRequiresBar{bar: bar}
}

func newBarRequiringFoo(foo fooBar) barFoo {
	return &barRequiresFoo{foo: foo}
}
