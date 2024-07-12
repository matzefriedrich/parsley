package pkg

import (
	"context"
	"github.com/matzefriedrich/parsley/internal"
	"github.com/matzefriedrich/parsley/pkg/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Resolver_ResolveWithOptions_inject_unregistered_service_instance(t *testing.T) {

	// Arrange
	sut := NewServiceRegistry()

	// Act
	_ = RegisterScoped(sut, NewFooConsumer)

	// Assert
	r := sut.BuildResolver()
	parsleyContext := internal.NewScopedContext(context.Background())
	consumer1, _ := r.ResolveWithOptions(parsleyContext, types.ServiceType[FooConsumer](), WithInstance[Foo](NewFoo()))
	assert.NotNil(t, consumer1)

	actual, ok := consumer1.(FooConsumer)
	assert.True(t, ok)
	assert.NotNil(t, actual)
}
