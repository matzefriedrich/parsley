package features

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MockFunction_String_returns_signature_if_present(t *testing.T) {
	// Arrange
	f := MockFunction{
		name:      "MyFunc",
		signature: "func(int) string",
	}

	// Act
	actual := f.String()

	// Assert
	assert.Equal(t, "func(int) string", actual)
}

func Test_MockFunction_String_returns_name_if_signature_is_missing(t *testing.T) {
	// Arrange
	f := MockFunction{
		name: "MyFunc",
	}

	// Act
	actual := f.String()

	// Assert
	assert.Equal(t, "MyFunc", actual)
}
