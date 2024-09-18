package generator

import (
	"github.com/matzefriedrich/parsley/internal/generator"
	"github.com/matzefriedrich/parsley/internal/reflection"
	"github.com/matzefriedrich/parsley/internal/tests/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_RegisterTypeModelFunctions_ensure_expected_methods(t *testing.T) {

	// Arrange
	expectedFunctionRegistrations := []string{
		"FormatType",
		"FormattedCallParameters",
		"FormattedParameters",
		"FormattedResultParameters",
		"FormattedResultTypes",
		"HasResults",
		"Signature",
	}

	registeredFunctions := make(map[string]struct{})
	target := mocks.NewGenericCodeGeneratorMock()
	target.AddTemplateFuncFunc = func(functions ...generator.TemplateFunction) error {
		for _, function := range functions {
			registeredFunctions[function.Name] = struct{}{}
		}
		return nil
	}

	// Act
	err := generator.RegisterTypeModelFunctions(target)

	// Assert
	assert.NoError(t, err)

	for _, name := range expectedFunctionRegistrations {
		_, registered := registeredFunctions[name]
		assert.True(t, registered)
	}
}

func Test_FormatType_scalar_value_type(t *testing.T) {
	// Arrange
	p := reflection.Parameter{Name: "p", Type: &reflection.ParameterType{Name: "string"}}
	// Act
	actual := generator.FormatType(p)
	// Assert
	assert.Equal(t, "string", actual)
}

func Test_FormatType_scalar_array_type(t *testing.T) {
	// Arrange
	p := reflection.Parameter{
		Name: "p",
		Type: &reflection.ParameterType{
			IsArray: true,
			Next: &reflection.ParameterType{
				Name: "string",
			},
		},
	}
	// Act
	actual := generator.FormatType(p)
	// Assert
	assert.Equal(t, "[]string", actual)
}

func Test_FormatType_array_pointer_type(t *testing.T) {
	// Arrange
	p := reflection.Parameter{
		Name: "p",
		Type: &reflection.ParameterType{
			IsArray: true,
			Next: &reflection.ParameterType{
				IsPointer: true,
				Next: &reflection.ParameterType{
					Name: "Arg",
				},
			},
		},
	}
	// Act
	actual := generator.FormatType(p)
	// Assert
	assert.Equal(t, "[]*Arg", actual)
}

func Test_FormattedCallParameters_single_parameter(t *testing.T) {
	// Arrange
	m := reflection.Method{
		Name: "SayHello",
		Parameters: []reflection.Parameter{
			{
				Name: "p", Type: &reflection.ParameterType{Name: "string"},
			},
		},
	}

	// Act
	actual := generator.FormattedCallParameters(m)
	// Assert
	assert.Equal(t, "p", actual)
}

func Test_FormattedCallParameters_multiple_parameter(t *testing.T) {
	// Arrange
	m := reflection.Method{
		Name: "SayHello",
		Parameters: []reflection.Parameter{
			{Name: "s", Type: &reflection.ParameterType{Name: "string"}},
			{Name: "b", Type: &reflection.ParameterType{Name: "bool"}},
		},
	}

	// Act
	actual := generator.FormattedCallParameters(m)
	// Assert
	assert.Equal(t, "s, b", actual)
}

func Test_FormattedParameters_single_parameter(t *testing.T) {
	// Arrange
	m := reflection.Method{
		Name: "SayHello",
		Parameters: []reflection.Parameter{
			{Name: "p", Type: &reflection.ParameterType{Name: "string"}},
		},
	}

	// Act
	actual := generator.FormattedParameters(m)
	// Assert
	assert.Equal(t, "p string", actual)
}

func Test_FormattedParameters_multiple_parameter(t *testing.T) {
	// Arrange
	m := reflection.Method{
		Name: "Method0",
		Parameters: []reflection.Parameter{
			{
				Name: "s",
				Type: &reflection.ParameterType{
					Name: "string",
				},
			},
			{
				Name: "buffer",
				Type: &reflection.ParameterType{
					IsArray: true,
					Next: &reflection.ParameterType{
						Name: "byte",
					},
				},
			},
			{
				Name: "b", Type: &reflection.ParameterType{
					Name: "bool"},
			},
		},
	}

	// Act
	actual := generator.FormattedParameters(m)
	// Assert
	assert.Equal(t, "s string, buffer []byte, b bool", actual)
}

func Test_FormattedResultParameters_single_parameter(t *testing.T) {
	// Arrange
	m := reflection.Method{
		Name: "SayHello",
		Parameters: []reflection.Parameter{
			{Name: "p", Type: &reflection.ParameterType{Name: "string"}},
		},
		Results: []reflection.Parameter{
			{Name: "result0", Type: &reflection.ParameterType{Name: "error"}},
		}}

	// Act
	actual := generator.FormattedResultParameters(m)
	// Assert
	assert.Equal(t, "result0", actual)
}

func Test_FormattedResultParameters_multiple_parameter(t *testing.T) {
	// Arrange
	m := reflection.Method{
		Name: "SayHello",
		Parameters: []reflection.Parameter{
			{Name: "s", Type: &reflection.ParameterType{Name: "string"}},
		},
		Results: []reflection.Parameter{
			{Name: "result0", Type: &reflection.ParameterType{Name: "string"}},
			{Name: "result1", Type: &reflection.ParameterType{Name: "error"}},
		},
	}
	// Act
	actual := generator.FormattedResultParameters(m)
	// Assert
	assert.Equal(t, "result0, result1", actual)
}
