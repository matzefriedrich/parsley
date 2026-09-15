package registration

import (
	"testing"

	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/types"
	"github.com/stretchr/testify/assert"
)

type groupableService struct{}

func newGroupableService() *groupableService {
	return &groupableService{}
}

func Test_RegisterSingleton_without_options_uses_the_default_group(t *testing.T) {
	// Arrange
	registry := registration.NewServiceRegistry()

	// Act
	err := registration.RegisterSingleton(registry, newGroupableService)

	// Assert
	assert.NoError(t, err)
	registrations, _ := registry.GetServiceRegistrations()
	assert.Equal(t, 1, len(registrations))
	assert.Equal(t, types.DefaultLifecycleGroup, registrations[0].LifecycleGroup())
}

func Test_CreateServiceRegistrationWithOptions_multiple_groups_returns_error(t *testing.T) {
	// Act
	_, err := registration.CreateServiceRegistrationWithOptions(newGroupableService, types.LifetimeSingleton,
		types.InLifecycleGroup("a"), types.InLifecycleGroup("b"))

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "multiple lifecycle groups")
}

func Test_RegisterSingletonWithOptions_assigns_the_lifecycle_group(t *testing.T) {
	// Arrange
	registry := registration.NewServiceRegistry()

	// Act
	err := registration.RegisterSingletonWithOptions(registry, newGroupableService, types.InLifecycleGroup("infrastructure"))

	// Assert
	assert.NoError(t, err)
	registrations, _ := registry.GetServiceRegistrations()
	assert.Equal(t, 1, len(registrations))
	assert.Equal(t, "infrastructure", registrations[0].LifecycleGroup())
}

func Test_RegisterSingletonWithOptions_empty_group_name_returns_error(t *testing.T) {
	// Arrange
	registry := registration.NewServiceRegistry()

	// Act
	err := registration.RegisterSingletonWithOptions(registry, newGroupableService, types.InLifecycleGroup(""))

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must not be empty")
}
