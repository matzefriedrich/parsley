package generator

import (
	"github.com/matzefriedrich/parsley/internal/generator"
	"github.com/matzefriedrich/parsley/internal/reflection"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_NewTemplateModelBuilder_from_empty_source_file_returns_error(t *testing.T) {

	// Arrange
	source := []byte("")
	sut := generator.NewTemplateModelBuilder(reflection.AstFromSource(source))

	// Act
	actual, err := sut.Build()

	// Assert
	assert.Error(t, err)
	assert.Nil(t, actual)
}

func Test_NewTemplateModelBuilder_from_minimal_source_file(t *testing.T) {

	// Arrange
	source := []byte("package main")
	sut := generator.NewTemplateModelBuilder(reflection.AstFromSource(source))

	// Act
	actual, err := sut.Build()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, actual)
}

func Test_NewTemplateModelBuilder_Build_multiple_interface_definitions(t *testing.T) {

	// Arrange
	source := []byte("package types\n" + "\n" +
		"type Service0 interface {\n" + "	Method0(s string)\n" + "}\n" + "\n" +
		"type Service1 interface {\n" + "	Method1() string\n" + "	Method2(data []bytes) (string, error)\n" +
		"}\n")

	accessor := reflection.AstFromSource(source)
	sut := generator.NewTemplateModelBuilder(accessor)

	// Act
	actual, err := sut.Build()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, actual)

	serviceInterface0 := actual.Interfaces[0]
	assert.Equal(t, "Service0", serviceInterface0.Name)

	serviceInterface1 := actual.Interfaces[1]
	assert.Equal(t, "Service1", serviceInterface1.Name)

	method1 := serviceInterface1.Methods[0]
	assert.Equal(t, "Method1", method1.Name)

	method2 := serviceInterface1.Methods[1]
	assert.Equal(t, "Method2", method2.Name)
}

func Test_NewTemplateModelBuilder_Build_collect_imports(t *testing.T) {

	// Arrange
	source := []byte("package types\n" + "\n" +
		"\n" +
		"import (\n" +
		"	\"fmt\"\n" +
		")")

	accessor := reflection.AstFromSource(source)
	sut := generator.NewTemplateModelBuilder(accessor)

	// Act
	actual, err := sut.Build()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, actual)

	assert.Equal(t, "\"fmt\"", actual.Imports[0])
}

func Test_NewTemplateModelBuilder_Build_collect_comments(t *testing.T) {

	// Arrange
	source := []byte("package types\n" + "\n" +
		"//parsley:ignore\n" +
		"type Greeter interface {\n" +
		"	SayHello(name string)\n" +
		"}")

	accessor := reflection.AstFromSource(source)
	sut := generator.NewTemplateModelBuilder(accessor)

	// Act
	actual, err := sut.Build()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, actual)

	assert.Equal(t, 1, len(actual.Comments))
	assert.Equal(t, "//parsley:ignore", actual.Comments[0].Text)
}

func Test_NewTemplateModelBuilder_Build_collect_struct_types(t *testing.T) {

	// Arrange
	source := []byte("package types\n" + "\n" +
		"type Greeter struct {\n" +
		"}")

	accessor := reflection.AstFromSource(source)
	sut := generator.NewTemplateModelBuilder(accessor)

	// Act
	actual, err := sut.Build()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, actual)

}
