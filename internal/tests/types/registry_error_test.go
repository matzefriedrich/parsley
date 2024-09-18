package types

import (
	"errors"
	"github.com/matzefriedrich/parsley/pkg/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_NewRegistryError_initializer_gets_invoked_with_expected_error_instance(t *testing.T) {

	// Arrange
	const expectedErrorMsg = "something went wrong"

	// Act
	actual := types.NewRegistryError(expectedErrorMsg, types.ForServiceType("Foo"))

	// Assert
	assert.ErrorIs(t, actual, &types.RegistryError{ParsleyError: types.ParsleyError{Msg: expectedErrorMsg}})

	var registryErr *types.RegistryError

	isRegistryErr := errors.As(actual, &registryErr)
	assert.True(t, isRegistryErr)

	matchesServiceType := registryErr.MatchesServiceType("Foo")
	assert.True(t, matchesServiceType)
}
