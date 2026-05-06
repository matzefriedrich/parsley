package registration

import (
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
	parsleyContext := resolving.NewScopedContext(t.Context())
	resolved, _ := resolving.ResolveRequiredService[FooConsumer](parsleyContext, r)
	assert.NotNil(t, resolved)
}

func Test_Registry_NewResolver_resolve_scoped_from_same_context_must_be_return_same_instance(t *testing.T) {

	// Arrange
	sut := registration.NewServiceRegistry()

	// Act
	_ = registration.RegisterSingleton(sut, newFoo)
	_ = registration.RegisterScoped(sut, newFooConsumer)

	// Assert
	r := resolving.NewResolver(sut)
	parsleyContext := resolving.NewScopedContext(t.Context())
	consumer1, _ := resolving.ResolveRequiredService[FooConsumer](parsleyContext, r)
	assert.NotNil(t, consumer1)

	consumer2, _ := resolving.ResolveRequiredService[FooConsumer](parsleyContext, r)
	assert.NotNil(t, consumer2)

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

	ctx := t.Context()

	parsleyContext1 := resolving.NewScopedContext(ctx)
	consumer1, _ := resolving.ResolveRequiredService[FooConsumer](parsleyContext1, r)

	assert.NotNil(t, consumer1)

	parsleyContext2 := resolving.NewScopedContext(ctx)
	consumer2, _ := resolving.ResolveRequiredService[FooConsumer](parsleyContext2, r)
	assert.NotNil(t, consumer2)

	assert.NotNil(t, consumer1)
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

func Test_Registry_RegisterModuleIf_registers_module_if_condition_is_true(t *testing.T) {

	// Arrange
	sut := registration.NewServiceRegistry()

	// Act
	fooModule := func(r types.ServiceRegistry) error {
		_ = registration.RegisterSingleton(r, newFoo)
		return nil
	}

	_ = sut.RegisterModuleIf(true, fooModule)

	fooRegistered := sut.IsRegistered(types.MakeServiceType[Foo]())

	// Assert
	assert.True(t, fooRegistered)
}

func Test_Registry_RegisterModuleIf_does_not_register_module_if_condition_is_false(t *testing.T) {

	// Arrange
	sut := registration.NewServiceRegistry()

	// Act
	fooModule := func(r types.ServiceRegistry) error {
		_ = registration.RegisterSingleton(r, newFoo)
		return nil
	}

	_ = sut.RegisterModuleIf(false, fooModule)

	fooRegistered := sut.IsRegistered(types.MakeServiceType[Foo]())

	// Assert
	assert.False(t, fooRegistered)
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

	actual, _ := resolving.ResolveRequiredService[Foo](t.Context(), r)

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
