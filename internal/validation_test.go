package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_IsNil_returns_false(t *testing.T) {

	// Arrange
	instance := &struct {
	}{}

	// Act
	actual := IsNil(instance)

	// Assert
	assert.Equal(t, false, actual)
}

func Test_IsNil_returns_false_for_empty_struct(t *testing.T) {

	// Arrange
	instance := struct {
	}{}

	// Act
	actual := IsNil(instance)

	// Assert
	assert.Equal(t, false, actual)
}

func Test_IsNil_returns_true(t *testing.T) {

	// Arrange
	var invalidInstance interface{}

	// Act
	actual := IsNil(invalidInstance)

	// Assert
	assert.Equal(t, true, actual)
}
