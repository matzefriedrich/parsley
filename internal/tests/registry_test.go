package tests

import (
	"context"
	"reflect"
	"testing"

	"github.com/matzefriedrich/parsley/internal/core"
	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/resolving"

	"github.com/matzefriedrich/parsley/pkg/types"

	"github.com/stretchr/testify/assert"
)

func Test_ServiceRegistry_register_types_with_different_lifetime_behavior(t *testing.T) {

	// Arrange
	sut := registration.NewServiceRegistry()

	// Act
	_ = registration.RegisterSingleton(sut, NewFoo)
	_ = registration.RegisterTransient(sut, NewFooConsumer)

	fooRegistered := sut.IsRegistered(registration.ServiceType[Foo]())
	fooConsumerRegistered := sut.IsRegistered(registration.ServiceType[FooConsumer]())

	// Assert
	assert.True(t, fooRegistered)
	assert.True(t, fooConsumerRegistered)
}

func Test_Registry_NewResolver_resolve_type_with_dependencies(t *testing.T) {

	// Arrange
	sut := registration.NewServiceRegistry()

	// Act
	_ = registration.RegisterTransient(sut, NewFoo)
	_ = registration.RegisterTransient(sut, NewFooConsumer)

	// Assert
	r := resolving.NewResolver(sut)
	parsleyContext := core.NewScopedContext(context.Background())
	resolved, _ := r.Resolve(parsleyContext, registration.ServiceType[FooConsumer]())
	assert.NotNil(t, resolved)

	actual, ok := resolved.(FooConsumer)
	assert.True(t, ok)
	assert.NotNil(t, actual)
}

func Test_Registry_NewResolver_resolve_scoped_from_same_context_must_be_return_same_instance(t *testing.T) {

	// Arrange
	sut := registration.NewServiceRegistry()

	// Act
	_ = registration.RegisterSingleton(sut, NewFoo)
	_ = registration.RegisterScoped(sut, NewFooConsumer)

	// Assert
	r := resolving.NewResolver(sut)
	parsleyContext := core.NewScopedContext(context.Background())
	consumer1, _ := r.Resolve(parsleyContext, registration.ServiceType[FooConsumer]())
	assert.NotNil(t, consumer1)

	consumer2, _ := r.Resolve(parsleyContext, registration.ServiceType[FooConsumer]())
	assert.NotNil(t, consumer2)

	actual, ok := consumer1.(FooConsumer)
	assert.True(t, ok)
	assert.NotNil(t, actual)
	assert.Equal(t, reflect.ValueOf(consumer1).Pointer(), reflect.ValueOf(consumer2).Pointer())
}

func Test_Registry_NewResolver_resolve_scoped_from_different_context_must_be_return_different_instance(t *testing.T) {

	// Arrange
	sut := registration.NewServiceRegistry()

	// Act
	_ = registration.RegisterSingleton(sut, NewFoo)
	_ = registration.RegisterScoped(sut, NewFooConsumer)

	// Assert
	r := resolving.NewResolver(sut)

	parsleyContext1 := core.NewScopedContext(context.Background())
	consumer1, _ := r.Resolve(parsleyContext1, registration.ServiceType[FooConsumer]())
	assert.NotNil(t, consumer1)

	parsleyContext2 := core.NewScopedContext(context.Background())
	consumer2, _ := r.Resolve(parsleyContext2, registration.ServiceType[FooConsumer]())
	assert.NotNil(t, consumer2)

	actual, ok := consumer1.(FooConsumer)
	assert.True(t, ok)
	assert.NotNil(t, actual)
	assert.NotEqual(t, reflect.ValueOf(consumer1).Pointer(), reflect.ValueOf(consumer2).Pointer())
}

func Test_Registry_RegisterModule_registers_collection_of_services(t *testing.T) {

	// Arrange
	sut := registration.NewServiceRegistry()

	// Act
	fooModule := func(r types.ServiceRegistry) error {
		_ = registration.RegisterSingleton(r, NewFoo)
		_ = registration.RegisterScoped(r, NewFooConsumer)
		return nil
	}

	_ = sut.RegisterModule(fooModule)

	fooRegistered := sut.IsRegistered(registration.ServiceType[Foo]())
	fooConsumerRegistered := sut.IsRegistered(registration.ServiceType[FooConsumer]())

	// Assert
	assert.True(t, fooRegistered)
	assert.True(t, fooConsumerRegistered)
}

func Test_Registry_RegisterInstance_registers_singleton_service_from_object(t *testing.T) {

	// Arrange
	sut := registration.NewServiceRegistry()

	instance := NewFoo()

	// Act
	_ = registration.RegisterInstance(sut, instance)

	fooRegistered := sut.IsRegistered(registration.ServiceType[Foo]())

	r := resolving.NewResolver(sut)

	// Arrange
	assert.True(t, fooRegistered)

	resolved, _ := r.Resolve(context.Background(), registration.ServiceType[Foo]())
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
