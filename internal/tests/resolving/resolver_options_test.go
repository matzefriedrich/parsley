package resolving

import (
	"testing"

	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/resolving"
	"github.com/matzefriedrich/parsley/pkg/types"
	"github.com/stretchr/testify/assert"
)

func Test_Resolver_ResolveWithOptions_inject_unregistered_service_instance(t *testing.T) {

	// Arrange
	sut := registration.NewServiceRegistry()

	// Act
	_ = registration.RegisterScoped(sut, newFooConsumer)

	// Assert
	r := resolving.NewResolver(sut)

	parsleyContext := resolving.NewScopedContext(t.Context())
	consumers, _ := r.ResolveWithOptions(parsleyContext, types.MakeServiceType[fooConsumer](), resolving.WithInstance[foo](newFoo()))
	assert.Equal(t, 1, len(consumers))

	consumer1 := consumers[0]
	assert.NotNil(t, consumer1)

	actual, ok := consumer1.(fooConsumer)
	assert.True(t, ok)
	assert.NotNil(t, actual)
}

type foo interface {
	Bar()
}

type fooService struct{}

func (f *fooService) Bar() {}

func newFoo() foo {
	return &fooService{}
}

type fooConsumer interface {
	FooBar()
}

type fooConsumerService struct {
	foo foo
}

func (fb *fooConsumerService) FooBar() {
	fb.foo.Bar()
}

func newFooConsumer(foo foo) fooConsumer {
	return &fooConsumerService{
		foo: foo,
	}
}
