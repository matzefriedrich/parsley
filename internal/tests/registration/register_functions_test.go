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

type registerFunc func(activatorFunc any, scope types.LifetimeScope) error

type registryMock struct {
	RegisterFunc registerFunc
}

func (r *registryMock) Register(activatorFunc any, scope types.LifetimeScope) error {
	return r.RegisterFunc(activatorFunc, scope)
}

var _ registration.SupportsRegisterActivatorFunc = (*registryMock)(nil)
