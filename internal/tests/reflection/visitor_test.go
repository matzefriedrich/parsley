package reflection

import (
	"testing"

	reflection2 "github.com/matzefriedrich/parsley/internal/reflection"
	"github.com/stretchr/testify/assert"
)

func Test_FileWalker_WalkSyntaxTree_build_Model_collect_interfaces(t *testing.T) {

	// Arrange
	fileVisitor := reflection2.NewFileVisitor()
	sut := reflection2.NewSyntaxWalker(fileVisitor)

	source := "" +
		"package main\n\n" +
		"" +
		"type SayHelloFunc func(name string)\n\n" +
		"" +
		"type Greeter interface {\n" +
		"	SayHello(name string)\n" +
		"}"

	fileAccessor := reflection2.AstFromSource([]byte(source))
	file, _ := fileAccessor()

	// Act
	err := sut.WalkSyntaxTree(file.File)

	// Assert
	assert.NoError(t, err)

	model, modelErr := fileVisitor.Model()
	assert.NoError(t, modelErr)
	assert.NotNil(t, model)

	assert.Equal(t, "Greeter", model.Interfaces[0].Name)
}

func Test_FileWalker_WalkSyntaxTree_build_Model_collect_func_types(t *testing.T) {

	// Arrange
	fileVisitor := reflection2.NewFileVisitor()
	sut := reflection2.NewSyntaxWalker(fileVisitor)

	source := "" +
		"package main\n\n" +
		"" +
		"type SayHelloFunc func(name string)\n\n"

	fileAccessor := reflection2.AstFromSource([]byte(source))
	file, _ := fileAccessor()

	// Act
	err := sut.WalkSyntaxTree(file.File)

	// Assert
	assert.NoError(t, err)

	model, modelErr := fileVisitor.Model()
	assert.NoError(t, modelErr)
	assert.NotNil(t, model)

	assert.Equal(t, "SayHelloFunc", model.FuncTypes[0].Name)
}

func Test_FileWalker_WalkSyntaxTree_build_Model_collect_struct_types(t *testing.T) {

	// Arrange
	fileVisitor := reflection2.NewFileVisitor()
	sut := reflection2.NewSyntaxWalker(fileVisitor)

	source := "" +
		"package main\n\n" +
		"type MyStruct struct {\n" +
		"	Name string\n" +
		"}\n"

	fileAccessor := reflection2.AstFromSource([]byte(source))
	file, _ := fileAccessor()

	// Act
	err := sut.WalkSyntaxTree(file.File)

	// Assert
	assert.NoError(t, err)

	model, modelErr := fileVisitor.Model()
	assert.NoError(t, modelErr)
	assert.NotNil(t, model)
}

func Test_FileWalker_WalkSyntaxTree_build_Model_collect_complex_interface(t *testing.T) {

	// Arrange
	fileVisitor := reflection2.NewFileVisitor()
	sut := reflection2.NewSyntaxWalker(fileVisitor)

	source := "" +
		"package main\n\n" +
		"type ComplexInterface interface {\n" +
		"	Variadic(args ...string)\n" +
		"	InterfaceParam(i interface{})\n" +
		"	NamedResults(a int) (err error, count int)\n" +
		"}\n"

	fileAccessor := reflection2.AstFromSource([]byte(source))
	file, _ := fileAccessor()

	// Act
	err := sut.WalkSyntaxTree(file.File)

	// Assert
	assert.NoError(t, err)

	model, modelErr := fileVisitor.Model()
	assert.NoError(t, modelErr)
	assert.NotNil(t, model)

	methods := model.Interfaces[0].Methods
	assert.Equal(t, 3, len(methods))

	// Variadic
	variadicMethod := methods[0]
	assert.True(t, variadicMethod.Parameters[0].Type.IsEllipsis)

	// InterfaceParam
	interfaceParamMethod := methods[1]
	assert.True(t, interfaceParamMethod.Parameters[0].Type.IsInterface)

	// NamedResults
	namedResultsMethod := methods[2]
	assert.Equal(t, "err", namedResultsMethod.Results[0].Name)
	assert.Equal(t, "count", namedResultsMethod.Results[1].Name)
}
