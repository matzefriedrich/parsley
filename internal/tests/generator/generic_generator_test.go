package generator

import (
	"fmt"
	"github.com/matzefriedrich/parsley/internal/generator"
	"github.com/matzefriedrich/parsley/internal/tests/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GenericGenerator_Generate_load_requested_template_via_loader_and_generates_output(t *testing.T) {
	// Arrange
	sut := generator.NewGenericCodeGenerator(func(name string) (string, error) {
		switch name {
		case "template":
			return "{{ .Msg }}", nil
		}
		return "", fmt.Errorf("template not found")
	})

	target := mocks.NewMemoryFile()

	// Act
	err := sut.Generate("template", struct{ Msg string }{Msg: "Hello"}, target)

	// Assert
	assert.NoError(t, err)

	actual := target.String()
	assert.Equal(t, "Hello", actual)
}
