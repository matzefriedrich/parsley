package types

import (
	"errors"
	"fmt"
	"github.com/matzefriedrich/parsley/pkg/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ParsleyError_errors_Is(t *testing.T) {

	// Arrange
	expectedErrorMsg := "something went wrong"
	sut := &types.ParsleyError{Msg: expectedErrorMsg}

	// Act
	actualError := sut.Error()

	// Assert
	assert.Equal(t, expectedErrorMsg, actualError)
	assert.ErrorIs(t, sut, errors.New(expectedErrorMsg))
}

func Test_ParsleyError_WithCause_Unwrap_returns_wrapped_error(t *testing.T) {

	// Arrange
	expectedErrorMsg := "something went wrong"
	sut := &types.ParsleyError{Msg: expectedErrorMsg}

	cause := fmt.Errorf("oops")
	types.WithCause(cause)(sut)

	// Act
	actualError := sut.Error()

	// Assert
	assert.Equal(t, expectedErrorMsg, actualError)
	assert.ErrorIs(t, sut, errors.New(expectedErrorMsg))
	assert.Equal(t, cause, sut.Unwrap())
}

func Test_ParsleyAggregateError_WithAggregatedCause_errors_Is_matches_any_of_the_wrapped_errors(t *testing.T) {
	// Arrange
	expectedErrorMsg := "something went wrong"
	sut := &types.ParsleyError{Msg: expectedErrorMsg}
	types.WithAggregatedCause(&types.ParsleyError{Msg: "oops"}, &types.ParsleyError{Msg: "ouch"})(sut)

	// Act
	actualError := sut.Error()

	// Assert
	assert.Equal(t, expectedErrorMsg, actualError)
	assert.ErrorIs(t, sut, &types.ParsleyAggregateError{Msg: expectedErrorMsg})
	assert.ErrorIs(t, sut, &types.ParsleyError{Msg: "oops"})
	assert.ErrorIs(t, sut, &types.ParsleyError{Msg: "ouch"})
}
