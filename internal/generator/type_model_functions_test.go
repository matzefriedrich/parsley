package generator

import (
	"github.com/matzefriedrich/parsley/internal/reflection"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_FormatType_formats_parameter_type(t *testing.T) {
	// Arrange
	p := reflection.Parameter{IsArray: false, TypeName: "byte", Name: "data"}

	// Act
	actual := FormatType(p)

	// Assert
	assert.Equal(t, "byte", actual)
}

func Test_FormatType_formats_parameter_array_type(t *testing.T) {
	// Arrange
	p := reflection.Parameter{IsArray: true, TypeName: "byte", Name: "data"}

	// Act
	actual := FormatType(p)

	// Assert
	assert.Equal(t, "[]byte", actual)
}

func Test_FormattedParameters_formats_parameters(t *testing.T) {
	// Arrange
	p := reflection.Parameter{IsArray: true, TypeName: "byte", Name: "data"}

	method := reflection.Method{
		Name:       "Method0",
		Parameters: []reflection.Parameter{p, {Name: "msg", TypeName: "string", IsArray: false}},
		Results:    []reflection.Parameter{},
	}

	// Act
	actual := FormattedParameters(method)

	// Assert
	assert.Equal(t, "data []byte, msg string", actual)
}
