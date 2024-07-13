package tests

import (
	"context"
	"github.com/matzefriedrich/parsley/internal/core"
	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/resolving"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Resolver_ResolveWithOptions_inject_unregistered_service_instance(t *testing.T) {

	// Arrange
	sut := registration.NewServiceRegistry()

	// Act
	_ = registration.RegisterScoped(sut, NewFooConsumer)

	// Assert
	r := resolving.NewResolver(sut)

	parsleyContext := core.NewScopedContext(context.Background())
	consumer1, _ := r.ResolveWithOptions(parsleyContext, registration.ServiceType[FooConsumer](), resolving.WithInstance[Foo](NewFoo()))
	assert.NotNil(t, consumer1)

	actual, ok := consumer1.(FooConsumer)
	assert.True(t, ok)
	assert.NotNil(t, actual)
}
