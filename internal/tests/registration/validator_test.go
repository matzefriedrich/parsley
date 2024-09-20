package registration

import (
	"errors"
	"fmt"
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

	var aggregateErr *types.ParsleyAggregateError
	if errors.As(err, &aggregateErr) {
		for _, e := range aggregateErr.Errors() {
			fmt.Println(e.Error())
		}
	}
}

func Test_Validator_Validate_detects_circular_dependencies_and_returns_error(t *testing.T) {

	// Arrange
	registry := registration.NewServiceRegistry()
	_ = registry.Register(newTestService2, types.LifetimeTransient)
	_ = registry.Register(newTestService3, types.LifetimeTransient)

	sut := registration.NewServiceRegistrationsValidator()

	// Act
	err := sut.Validate(registry)

	// Assert
	assert.ErrorIs(t, err, registration.ErrCircularDependencyDetected)

	var aggregateErr *types.ParsleyAggregateError
	if errors.As(err, &aggregateErr) {
		for _, e := range aggregateErr.Errors() {
			fmt.Println(e.Error())
		}
	}
}

type testService struct {
	greeter features.Greeter
}

func newTestService(greeter features.Greeter) *testService {
	return &testService{
		greeter: greeter,
	}
}

type testService2 struct {
	other *testService3
}

func newTestService2(other *testService3) *testService2 {
	return &testService2{
		other: other,
	}
}

type testService3 struct {
	other *testService2
}

func newTestService3(other *testService2) *testService3 {
	return &testService3{
		other: other,
	}
}
