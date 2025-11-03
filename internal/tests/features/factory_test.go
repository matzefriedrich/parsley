package features

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/matzefriedrich/parsley/pkg/features"
	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/resolving"
	"github.com/matzefriedrich/parsley/pkg/types"
	"github.com/stretchr/testify/assert"
)

func Test_register_factory_for_type_returns_transient_instance(t *testing.T) {

	// Arrange
	registry := registration.NewServiceRegistry()
	registry.Register(newGreeterWithState, types.LifetimeTransient)
	_ = features.RegisterFactory[Greeter](registry, types.LifetimeSingleton)

	resolver := resolving.NewResolver(registry)

	ctx := t.Context()

	// Act
	factory, _ := resolving.ResolveRequiredService[features.FactoryFunc[Greeter]](ctx, resolver)

	scopedContext := resolving.NewScopedContext(ctx)
	actual, _ := factory(scopedContext)
	other, _ := factory(scopedContext)

	// Assert
	assert.NotNil(t, actual)
	assert.NotNil(t, other)

	actualValue := reflect.ValueOf(actual)
	actualInstancePointer := actualValue.Pointer()

	otherValue := reflect.ValueOf(other)
	otherInstancePointer := otherValue.Pointer()

	assert.NotEqual(t, actualInstancePointer, otherInstancePointer)
}

func Test_register_factory_for_type_returns_same_instance_per_scope(t *testing.T) {

	// Arrange
	registry := registration.NewServiceRegistry()
	registry.Register(newGreeterWithState, types.LifetimeScoped)
	_ = features.RegisterFactory[Greeter](registry, types.LifetimeSingleton)

	resolver := resolving.NewResolver(registry)

	ctx := t.Context()

	// Act
	factory, _ := resolving.ResolveRequiredService[features.FactoryFunc[Greeter]](ctx, resolver)

	scopedContext := resolving.NewScopedContext(ctx)
	actual, _ := factory(scopedContext)
	other, _ := factory(scopedContext)

	// Assert
	assert.NotNil(t, actual)
	assert.NotNil(t, other)

	actualValue := reflect.ValueOf(actual)
	actualInstancePointer := actualValue.Pointer()

	otherValue := reflect.ValueOf(other)
	otherInstancePointer := otherValue.Pointer()

	assert.Equal(t, actualInstancePointer, otherInstancePointer)
}

type statefulGreeter struct {
	v int
}

func (g *statefulGreeter) SayNothing() {
}

func (g *statefulGreeter) SayHello(name string, _ bool) (string, error) {
	return "Hello " + name, nil
}

var _ Greeter = &greeter{}

func newGreeterWithState() Greeter {
	return &statefulGreeter{
		v: rand.Int(), // needed to work around the Go-compilers optimization (allocation elision and escape analysis)
	}
}
