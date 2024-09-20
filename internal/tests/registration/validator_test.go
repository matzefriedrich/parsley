package registration

import (
	"github.com/matzefriedrich/parsley/internal/tests/features"
	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Validator_Validate_on_empty_registry_does_not_return_error(t *testing.T) {

	// Arrange
	registry := registration.NewServiceRegistry()
	sut := registration.NewServiceRegistrationsValidator()

	// Act
	err := sut.Validate(registry)

	// Assert
	assert.NoError(t, err)
}

func Test_Validator_Validate_single_service_with_no_dependency_does_not_return_error(t *testing.T) {

	// Arrange
	registry := registration.NewServiceRegistry()
	_ = registry.Register(features.NewGreeterMock, types.LifetimeTransient)

	sut := registration.NewServiceRegistrationsValidator()

	// Act
	err := sut.Validate(registry)

	// Assert
	assert.NoError(t, err)
}

func Test_Validator_Validate_missing_dependency_does_return_error(t *testing.T) {

	// Arrange
	registry := registration.NewServiceRegistry()
	_ = registry.Register(newTestService, types.LifetimeTransient)

	sut := registration.NewServiceRegistrationsValidator()

	// Act
	err := sut.Validate(registry)

	// Assert
	assert.ErrorIs(t, err, registration.ErrRegistryMissesRequiredServiceRegistrations)
}

type testService struct {
	greeter features.Greeter
}

func newTestService(greeter features.Greeter) *testService {
	return &testService{
		greeter: greeter,
	}
}
