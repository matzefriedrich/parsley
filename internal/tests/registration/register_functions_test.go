package registration

import (
	"fmt"
	"github.com/matzefriedrich/parsley/internal/tests/features"
	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_RegisterTransient_returns_error(t *testing.T) {

	// Arrange
	registryMock := &registryMock{}
	registryMock.RegisterFunc = func(activatorFunc any, scope types.LifetimeScope) error {
		if scope == types.LifetimeTransient {
			return fmt.Errorf("oops")
		}
		return nil
	}

	// Act
	err := registration.RegisterTransient(registryMock, features.NewGreeterMock)

	// Assert
	assert.Error(t, err)
}

func Test_RegisterScoped_returns_error(t *testing.T) {

	// Arrange
	registryMock := &registryMock{}
	registryMock.RegisterFunc = func(activatorFunc any, scope types.LifetimeScope) error {
		if scope == types.LifetimeScoped {
			return fmt.Errorf("oops")
		}
		return nil
	}

	// Act
	err := registration.RegisterScoped(registryMock, features.NewGreeterMock)

	// Assert
	assert.Error(t, err)
}

func Test_RegisterSingleton_returns_error(t *testing.T) {

	// Arrange
	registryMock := &registryMock{}
	registryMock.RegisterFunc = func(activatorFunc any, scope types.LifetimeScope) error {
		if scope == types.LifetimeSingleton {
			return fmt.Errorf("oops")
		}
		return nil
	}

	// Act
	err := registration.RegisterSingleton(registryMock, features.NewGreeterMock)

	// Assert
	assert.Error(t, err)
}

func Test_RegisterSingletonWithOptions_forwards_scope_and_options(t *testing.T) {
	// Arrange
	registryMock := &registryMock{}
	var capturedScope types.LifetimeScope
	var capturedOptions []types.LifecycleOption
	registryMock.RegisterWithOptionsFunc = func(activatorFunc any, scope types.LifetimeScope, options ...types.LifecycleOption) error {
		capturedScope = scope
		capturedOptions = options
		return nil
	}
	option := types.InLifecycleGroup("transport")

	// Act
	err := registration.RegisterSingletonWithOptions(registryMock, features.NewGreeterMock, option)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, types.LifetimeSingleton, capturedScope)
	assert.Equal(t, 1, len(capturedOptions))
}

func Test_RegisterTransientWithOptions_forwards_scope_and_options(t *testing.T) {
	// Arrange
	registryMock := &registryMock{}
	var capturedScope types.LifetimeScope
	registryMock.RegisterWithOptionsFunc = func(activatorFunc any, scope types.LifetimeScope, options ...types.LifecycleOption) error {
		capturedScope = scope
		return nil
	}

	// Act
	err := registration.RegisterTransientWithOptions(registryMock, features.NewGreeterMock, types.InLifecycleGroup("infrastructure"))

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, types.LifetimeTransient, capturedScope)
}

func Test_RegisterScopedWithOptions_forwards_scope(t *testing.T) {
	// Arrange
	registryMock := &registryMock{}
	var capturedScope types.LifetimeScope
	registryMock.RegisterWithOptionsFunc = func(activatorFunc any, scope types.LifetimeScope, options ...types.LifecycleOption) error {
		capturedScope = scope
		return nil
	}

	// Act
	err := registration.RegisterScopedWithOptions(registryMock, features.NewGreeterMock, types.InLifecycleGroup("infrastructure"))

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, types.LifetimeScoped, capturedScope)
}

type registerFunc func(activatorFunc any, scope types.LifetimeScope) error

type registerWithOptionsFunc func(activatorFunc any, scope types.LifetimeScope, options ...types.LifecycleOption) error

type registryMock struct {
	RegisterFunc            registerFunc
	RegisterWithOptionsFunc registerWithOptionsFunc
}

func (r *registryMock) Register(activatorFunc any, scope types.LifetimeScope) error {
	return r.RegisterFunc(activatorFunc, scope)
}

func (r *registryMock) RegisterWithOptions(activatorFunc any, scope types.LifetimeScope, options ...types.LifecycleOption) error {
	return r.RegisterWithOptionsFunc(activatorFunc, scope, options...)
}

var _ registration.SupportsRegisterActivatorFunc = (*registryMock)(nil)
var _ registration.SupportsRegisterActivatorFuncWithOptions = (*registryMock)(nil)
