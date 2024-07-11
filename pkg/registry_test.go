package pkg

import (
	"context"
	"reflect"
	"testing"

	"github.com/matzefriedrich/parsley/internal"

	"github.com/matzefriedrich/parsley/pkg/types"

	"github.com/stretchr/testify/assert"
)

func Test_ServiceRegistry_register_types(t *testing.T) {

	// Arrange
	sut := NewServiceRegistry()

	// Act
	_ = RegisterSingleton(sut, NewFoo)
	_ = RegisterTransient(sut, NewFooConsumer)

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
	_ = RegisterTransient(sut, NewFoo)
	_ = RegisterTransient(sut, NewFooConsumer)

	// Assert
	r := sut.BuildResolver()
	parsleyContext := internal.NewScopedContext(context.Background())
	resolved, _ := r.Resolve(parsleyContext, types.ServiceType[FooConsumer]())
	assert.NotNil(t, resolved)

	actual, ok := resolved.(FooConsumer)
	assert.True(t, ok)
	assert.NotNil(t, actual)
}

func Test_Registry_BuildResolver_resolve_scoped_from_same_context_must_be_return_same_instance(t *testing.T) {

	// Arrange
	sut := NewServiceRegistry()

	// Act
	_ = RegisterSingleton(sut, NewFoo)
	_ = RegisterScoped(sut, NewFooConsumer)

	// Assert
	r := sut.BuildResolver()
	parsleyContext := internal.NewScopedContext(context.Background())
	consumer1, _ := r.Resolve(parsleyContext, types.ServiceType[FooConsumer]())
	assert.NotNil(t, consumer1)

	consumer2, _ := r.Resolve(parsleyContext, types.ServiceType[FooConsumer]())
	assert.NotNil(t, consumer2)

	actual, ok := consumer1.(FooConsumer)
	assert.True(t, ok)
	assert.NotNil(t, actual)
	assert.Equal(t, reflect.ValueOf(consumer1).Pointer(), reflect.ValueOf(consumer2).Pointer())
}

func Test_Registry_BuildResolver_resolve_scoped_from_different_context_must_be_return_different_instance(t *testing.T) {

	// Arrange
	sut := NewServiceRegistry()

	// Act
	_ = RegisterSingleton(sut, NewFoo)
	_ = RegisterScoped(sut, NewFooConsumer)

	// Assert
	r := sut.BuildResolver()
	parsleyContext1 := internal.NewScopedContext(context.Background())
	consumer1, _ := r.Resolve(parsleyContext1, types.ServiceType[FooConsumer]())
	assert.NotNil(t, consumer1)

	parsleyContext2 := internal.NewScopedContext(context.Background())
	consumer2, _ := r.Resolve(parsleyContext2, types.ServiceType[FooConsumer]())
	assert.NotNil(t, consumer2)

	actual, ok := consumer1.(FooConsumer)
	assert.True(t, ok)
	assert.NotNil(t, actual)
	assert.NotEqual(t, reflect.ValueOf(consumer1).Pointer(), reflect.ValueOf(consumer2).Pointer())
}

func Test_Registry_RegisterModule_registers_collection_of_services(t *testing.T) {

	// Arrange
	sut := NewServiceRegistry()

	// Act
	fooModule := func(r types.ServiceRegistry) error {
		_ = RegisterSingleton(r, NewFoo)
		_ = RegisterScoped(r, NewFooConsumer)
		return nil
	}

	_ = sut.RegisterModule(fooModule)

	fooRegistered := sut.IsRegistered(types.ServiceType[Foo]())
	fooConsumerRegistered := sut.IsRegistered(types.ServiceType[FooConsumer]())

	// Assert
	assert.True(t, fooRegistered)
	assert.True(t, fooConsumerRegistered)
}

func Test_Registry_RegisterInstance_registers_object(t *testing.T) {

	// Arrange
	sut := NewServiceRegistry()

	instance := NewFoo()

	// Act
	_ = RegisterInstance(sut, instance)

	fooRegistered := sut.IsRegistered(types.ServiceType[Foo]())

	r := sut.BuildResolver()

	// Arrange
	assert.True(t, fooRegistered)

	resolved, _ := r.Resolve(context.Background(), types.ServiceType[Foo]())
	assert.NotNil(t, resolved)

	actual, ok := resolved.(Foo)
	assert.True(t, ok)
	assert.NotNil(t, actual)
	assert.Equal(t, reflect.ValueOf(instance).Pointer(), reflect.ValueOf(actual).Pointer())

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
