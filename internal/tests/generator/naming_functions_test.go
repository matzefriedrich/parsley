package generator

import (
	"github.com/matzefriedrich/parsley/internal/generator"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_RegisterNamingFunctions_(t *testing.T) {

	// Arrange
	expectedFunctionRegistrations := []string{"asPrivate", "asPublic"}

	registeredFunctions := make(map[string]struct{})
	target := &genericCodeGeneratorMock{
		AddTemplateFuncFunc: func(functions ...generator.TemplateFunction) error {
			for _, function := range functions {
				registeredFunctions[function.Name] = struct{}{}
			}
			return nil
		},
	}

	// Act
	err := generator.RegisterNamingFunctions(target)

	// Assert
	assert.NoError(t, err)

	for _, name := range expectedFunctionRegistrations {
		_, registered := registeredFunctions[name]
		assert.True(t, registered)
	}
}

func Test_MakePrivate_changes_PascalCase_to_camelCase(t *testing.T) {
	// Arrange
	const name = "SymbolName"
	// Act
	actual := generator.MakePrivate(name)
	// Assert
	assert.Equal(t, "symbolName", actual)
}

func Test_MakePrivate_does_not_change_camelCase(t *testing.T) {
	// Arrange
	const name = "symbolName"
	// Act
	actual := generator.MakePrivate(name)
	// Assert
	assert.Equal(t, "symbolName", actual)
}

func Test_MakePublic_changes_camelCase_to_PascalCase(t *testing.T) {
	// Arrange
	const name = "symbolName"
	// Act
	actual := generator.MakePublic(name)
	// Assert
	assert.Equal(t, "SymbolName", actual)
}

func Test_MakePublic_does_not_change_PascalCase(t *testing.T) {
	// Arrange
	const name = "SymbolName"
	// Act
	actual := generator.MakePublic(name)
	// Assert
	assert.Equal(t, "SymbolName", actual)
}
