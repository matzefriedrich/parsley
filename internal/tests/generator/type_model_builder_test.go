package generator

import (
	"github.com/matzefriedrich/parsley/internal/generator"
	"github.com/matzefriedrich/parsley/internal/reflection"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_AstFromSource_empty_source_file_returns_error(t *testing.T) {

	// Arrange
	source := []byte("")
	accessor := reflection.AstFromSource(source)

	// Act
	file, err := accessor()

	// Assert
	assert.Error(t, err)
	assert.Nil(t, file)
}

func Test_NewTemplateModelBuilder_from_minimal_source_file(t *testing.T) {

	// Arrange
	source := []byte("package main")

	accessor := reflection.AstFromSource(source)
	file, _ := accessor()

	sut := generator.NewTemplateModelBuilder(file.File)

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
	file, _ := accessor()

	sut := generator.NewTemplateModelBuilder(file.File)

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
	file, _ := accessor()

	sut := generator.NewTemplateModelBuilder(file.File)

	// Act
	actual, err := sut.Build()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, actual)

	assert.Equal(t, "fmt", actual.Imports[0])
}

func Test_NewTemplateModelBuilder_Build_collect_comments(t *testing.T) {

	// Arrange
	source := []byte("package types\n" + "\n" +
		"//parsley:ignore\n" +
		"type Greeter interface {\n" +
		"	SayHello(name string)\n" +
		"}")

	accessor := reflection.AstFromSource(source)
	file, _ := accessor()

	sut := generator.NewTemplateModelBuilder(file.File)

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
	file, _ := accessor()

	sut := generator.NewTemplateModelBuilder(file.File)

	// Act
	actual, err := sut.Build()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, actual)
}

func Test_NewTemplateModelBuilder_Build_interface_method_with_pointer_parameters(t *testing.T) {

	// Arrange
	source := []byte("package types\n" + "\n" +
		"import (" + "\n" +
		"	\"net/http\"" + "\n" +
		")" + "\n\n" +
		"type HttpClient interface {\n" +
		"	Do(req *http.Request) (*http.Response, error)" + "\n" +
		"}")

	accessor := reflection.AstFromSource(source)
	file, _ := accessor()

	sut := generator.NewTemplateModelBuilder(file.File)

	// Act
	actual, err := sut.Build()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, actual)

	httpClientInterface := actual.Interfaces[0]
	assert.Equal(t, "HttpClient", httpClientInterface.Name)

	doMethod := httpClientInterface.Methods[0]
	assert.Equal(t, "Do", doMethod.Name)

	reqParameter := doMethod.Parameters[0]
	assert.Equal(t, "req", reqParameter.Name)
	// assert.Equal(t, "http.Request", reqParameter.TypeName)
	// assert.True(t, reqParameter.IsPointer)

	// responseResult := doMethod.Results[0]
	// assert.Equal(t, "http.Response", responseResult.TypeName)
	// assert.True(t, responseResult.IsPointer)
}

func Test_NewTemplateModelBuilder_Build_interface_method_array_pointer_parameter(t *testing.T) {

	// Arrange
	source := []byte("package types\n" + "\n" +
		"import (" + "\n" +
		"	\"types\"" + "\n" +
		")" + "\n\n" +
		"type Service interface {\n" +
		"	Method0(args *[]*types.Arg)" + "\n" +
		"}")

	accessor := reflection.AstFromSource(source)
	file, _ := accessor()

	sut := generator.NewTemplateModelBuilder(file.File)

	// Act
	actual, err := sut.Build()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, actual)

	method := actual.Interfaces[0].Methods[0]
	signature := generator.Signature(method)
	assert.Equal(t, "Method0(args *[]*types.Arg)", signature)
}
