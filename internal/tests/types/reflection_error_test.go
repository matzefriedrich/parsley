package types

import (
	"errors"
	"github.com/matzefriedrich/parsley/pkg/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_NewReflectionError_initializer_gets_invoked_with_expected_error_instance(t *testing.T) {

	// Arrange
	const expectedErrorMsg = "something went wrong"
	var initializerInvokedWithExpectedError = false

	// Act
	actual := types.NewReflectionError(expectedErrorMsg, func(e error) {
		var funqErr *types.ParsleyError
		ok := errors.As(e, &funqErr)
		if ok {
			initializerInvokedWithExpectedError = true
		}
	})

	// Assert
	assert.ErrorIs(t, actual, &types.ParsleyError{Msg: expectedErrorMsg})
	assert.True(t, initializerInvokedWithExpectedError)
}
