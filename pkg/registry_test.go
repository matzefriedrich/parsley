package pkg

import (
	"context"
	"testing"

	"github.com/matzefriedrich/parsley/pkg/types"

	"github.com/stretchr/testify/assert"
)

func Test_ServiceRegistry_register_types(t *testing.T) {

	// Arrange
	sut := NewServiceRegistry()

	// Act
	_ = sut.Register(NewFoo)
	_ = sut.Register(NewFooConsumer)

	fooRegistered := sut.IsRegistered(types.ServiceType[Foo]())
	fooConsumerRegistered := sut.IsRegistered(types.ServiceType[FooConsumer]())

	// Assert
	assert.True(t, fooRegistered)
	assert.True(t, fooConsumerRegistered)
}

func Test_Registry_BuildResolver_resolve_type_with_dependencies(t *testing.T) {

	// Arrange
	sut := NewServiceRegistry()
	// Act
	_ = sut.Register(NewFoo)
	_ = sut.Register(NewFooConsumer)

	// Assert
	r := sut.BuildResolver()
	resolved, _ := r.Resolve(context.Background(), types.ServiceType[FooConsumer]())
	assert.NotNil(t, resolved)

	actual, ok := resolved.(FooConsumer)
	assert.True(t, ok)
	assert.NotNil(t, actual)
}

type Foo interface {
	Bar()
}

type foo struct{}

func (f *foo) Bar() {}

func NewFoo() Foo {
	return &foo{}
}

type FooConsumer interface {
	FooBar()
}

type fooConsumer struct {
	foo Foo
}

func (fb *fooConsumer) FooBar() {
	fb.foo.Bar()
}

func NewFooConsumer(foo Foo) FooConsumer {
	return &fooConsumer{
		foo: foo,
	}
}
