package registration

import (
	"context"
	"reflect"
	"testing"

	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/resolving"

	"github.com/matzefriedrich/parsley/pkg/types"

	"github.com/stretchr/testify/assert"
)

func Test_ServiceRegistry_register_types_with_different_lifetime_behavior(t *testing.T) {

	// Arrange
	sut := registration.NewServiceRegistry()

	// Act
	_ = registration.RegisterSingleton(sut, newFoo)
	_ = registration.RegisterTransient(sut, newFooConsumer)

	fooRegistered := sut.IsRegistered(types.MakeServiceType[Foo]())
	fooConsumerRegistered := sut.IsRegistered(types.MakeServiceType[FooConsumer]())

	// Assert
	assert.True(t, fooRegistered)
	assert.True(t, fooConsumerRegistered)
}

func Test_Registry_NewResolver_resolve_type_with_dependencies(t *testing.T) {

	// Arrange
	sut := registration.NewServiceRegistry()

	// Act
	_ = registration.RegisterTransient(sut, newFoo)
	_ = registration.RegisterTransient(sut, newFooConsumer)

	// Assert
	r := resolving.NewResolver(sut)
	parsleyContext := resolving.NewScopedContext(context.Background())
	resolved, _ := resolving.ResolveRequiredService[FooConsumer](r, parsleyContext)
	assert.NotNil(t, resolved)

	actual, ok := resolved.(FooConsumer)
	assert.True(t, ok)
	assert.NotNil(t, actual)
}

func Test_Registry_NewResolver_resolve_scoped_from_same_context_must_be_return_same_instance(t *testing.T) {

	// Arrange
	sut := registration.NewServiceRegistry()

	// Act
	_ = registration.RegisterSingleton(sut, newFoo)
	_ = registration.RegisterScoped(sut, newFooConsumer)

	// Assert
	r := resolving.NewResolver(sut)
	parsleyContext := resolving.NewScopedContext(context.Background())
	consumer1, _ := resolving.ResolveRequiredService[FooConsumer](r, parsleyContext)
	assert.NotNil(t, consumer1)

	consumer2, _ := resolving.ResolveRequiredService[FooConsumer](r, parsleyContext)
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
	_ = registration.RegisterSingleton(sut, newFoo)
	_ = registration.RegisterScoped(sut, newFooConsumer)

	// Assert
	r := resolving.NewResolver(sut)

	parsleyContext1 := resolving.NewScopedContext(context.Background())
	consumer1, _ := resolving.ResolveRequiredService[FooConsumer](r, parsleyContext1)

	assert.NotNil(t, consumer1)

	parsleyContext2 := resolving.NewScopedContext(context.Background())
	consumer2, _ := resolving.ResolveRequiredService[FooConsumer](r, parsleyContext2)
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
		_ = registration.RegisterSingleton(r, newFoo)
		_ = registration.RegisterScoped(r, newFooConsumer)
		return nil
	}

	_ = sut.RegisterModule(fooModule)

	fooRegistered := sut.IsRegistered(types.MakeServiceType[Foo]())
	fooConsumerRegistered := sut.IsRegistered(types.MakeServiceType[FooConsumer]())

	// Assert
	assert.True(t, fooRegistered)
	assert.True(t, fooConsumerRegistered)
}

func Test_Registry_RegisterInstance_registers_singleton_service_from_object(t *testing.T) {

	// Arrange
	sut := registration.NewServiceRegistry()

	instance := newFoo()

	// Act
	_ = registration.RegisterInstance(sut, instance)

	fooRegistered := sut.IsRegistered(types.MakeServiceType[Foo]())

	r := resolving.NewResolver(sut)

	// Arrange
	assert.True(t, fooRegistered)

	resolved, _ := resolving.ResolveRequiredService[Foo](r, context.Background())
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

func newFoo() Foo {
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

func newFooConsumer(foo Foo) FooConsumer {
	return &fooConsumer{
		foo: foo,
	}
}
