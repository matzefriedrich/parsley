package registration

import (
	"context"
	"errors"
	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/resolving"
	"github.com/matzefriedrich/parsley/pkg/types"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func Test_Registry_register_multiple_transient_types(t *testing.T) {

	// Arrange
	sut := registration.NewServiceRegistry()

	// Act
	_ = registration.RegisterTransient(sut, newFoo3)
	registerErr := registration.RegisterTransient(sut, newFoo4)
	if errors.Is(registerErr, types.ErrTypeAlreadyRegistered) {
		assert.FailNow(t, "type is already registered")
	}

	r := resolving.NewResolver(sut)

	// Assert
	resolvedServices, err := resolving.ResolveRequiredServices[multiFoo](r, context.Background())

	assert.NoError(t, err)
	assert.Len(t, resolvedServices, 2)

	var foo3Instance, foo4Instance multiFoo
	for _, service := range resolvedServices {
		switch service.(type) {
		case *foo3:
			foo3Instance = service.(*foo3)
		case *foo4:
			foo4Instance = service.(*foo4)
		}
	}

	assert.NotNil(t, foo3Instance)
	assert.Equal(t, "foo3", foo3Instance.Bar())

	assert.NotNil(t, foo4Instance)
	assert.Equal(t, "foo4", foo4Instance.Bar())
}

func Test_Registry_register_multiple_types_mixed_lifetime_scopes(t *testing.T) {

	// Arrange
	sut := registration.NewServiceRegistry()

	// Act
	_ = registration.RegisterTransient(sut, newFoo3)
	registerErr := registration.RegisterSingleton(sut, newFoo4)
	if errors.Is(registerErr, types.ErrTypeAlreadyRegistered) {
		assert.FailNow(t, "type is already registered")
	}

	r := resolving.NewResolver(sut)
	resolvedServices1, err := resolving.ResolveRequiredServices[multiFoo](r, context.Background())
	resolvedServices2, _ := resolving.ResolveRequiredServices[multiFoo](r, context.Background())

	// Assert
	assert.NoError(t, err)
	assert.Len(t, resolvedServices1, 2)

	var foo3Instance1, foo4Instance1 multiFoo
	for _, service := range resolvedServices1 {
		switch service.(type) {
		case *foo3:
			foo3Instance1 = service.(*foo3)
		case *foo4:
			foo4Instance1 = service.(*foo4)
		}
	}

	assert.NotNil(t, foo3Instance1)
	assert.Equal(t, "foo3", foo3Instance1.Bar())

	assert.NotNil(t, foo4Instance1)
	assert.Equal(t, "foo4", foo4Instance1.Bar())

	var foo3Instance2, foo4Instance2 multiFoo
	for _, service := range resolvedServices2 {
		switch service.(type) {
		case *foo3:
			foo3Instance2 = service.(*foo3)
		case *foo4:
			foo4Instance2 = service.(*foo4)
		}
	}

	assert.NotNil(t, foo3Instance2)
	assert.NotEqual(t, reflect.ValueOf(foo3Instance1).Pointer(), reflect.ValueOf(foo3Instance2).Pointer())

	assert.NotNil(t, foo4Instance2)
	assert.Equal(t, reflect.ValueOf(foo4Instance1).Pointer(), reflect.ValueOf(foo4Instance2).Pointer())
}

func Test_Registry_register_multiple_types_mixed_lifetime_scopes_2(t *testing.T) {

	// Arrange
	sut := registration.NewServiceRegistry()

	// Act
	_ = registration.RegisterTransient(sut, newFoo3)
	registerErr := registration.RegisterScoped(sut, newFoo4)
	if errors.Is(registerErr, types.ErrTypeAlreadyRegistered) {
		assert.FailNow(t, "type is already registered")
	}

	r := resolving.NewResolver(sut)
	scope := resolving.NewScopedContext(context.Background())
	resolvedServices1, err := resolving.ResolveRequiredServices[multiFoo](r, scope)
	resolvedServices2, _ := resolving.ResolveRequiredServices[multiFoo](r, scope)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, resolvedServices1, 2)

	var foo3Instance1, foo4Instance1 multiFoo
	for _, service := range resolvedServices1 {
		switch service.(type) {
		case *foo3:
			foo3Instance1 = service.(*foo3)
		case *foo4:
			foo4Instance1 = service.(*foo4)
		}
	}

	assert.NotNil(t, foo3Instance1)
	assert.Equal(t, "foo3", foo3Instance1.Bar())

	assert.NotNil(t, foo4Instance1)
	assert.Equal(t, "foo4", foo4Instance1.Bar())

	var foo3Instance2, foo4Instance2 multiFoo
	for _, service := range resolvedServices2 {
		switch service.(type) {
		case *foo3:
			foo3Instance2 = service.(*foo3)
		case *foo4:
			foo4Instance2 = service.(*foo4)
		}
	}

	assert.NotNil(t, foo3Instance2)
	assert.NotEqual(t, reflect.ValueOf(foo3Instance1).Pointer(), reflect.ValueOf(foo3Instance2).Pointer())

	assert.NotNil(t, foo4Instance2)
	assert.Equal(t, reflect.ValueOf(foo4Instance1).Pointer(), reflect.ValueOf(foo4Instance2).Pointer())
}

type multiFoo interface {
	Bar() string
}

type foo3 struct {
	name string
}

func newFoo3() multiFoo {
	return &foo3{
		name: "foo3",
	}
}

func (f foo3) Bar() string {
	return f.name
}

type foo4 struct {
	name string
}

func newFoo4() multiFoo {
	return &foo4{
		name: "foo4",
	}
}
func (f foo4) Bar() string {
	return f.name
}

var _ multiFoo = &foo3{}
var _ multiFoo = &foo4{}
