package types

import (
	"errors"
	"testing"

	"github.com/matzefriedrich/parsley/pkg/types"
	"github.com/stretchr/testify/assert"
)

func Test_NewRegistryError_initializer_gets_invoked_with_expected_error_instance(t *testing.T) {

	// Arrange
	const expectedErrorMsg = "something went wrong"

	// Act
	actual := types.NewRegistryError(expectedErrorMsg, types.ForServiceTypeByName("Foo"))

	// Assert
	assert.ErrorIs(t, actual, &types.RegistryError{ParsleyError: types.ParsleyError{Msg: expectedErrorMsg}})

	var registryErr *types.RegistryError

	isRegistryErr := errors.As(actual, &registryErr)
	assert.True(t, isRegistryErr)

	matchesServiceType := registryErr.ServiceTypeName() == "Foo"
	assert.True(t, matchesServiceType)
}
