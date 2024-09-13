package reflection

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_FileWalker_WalkSyntaxTree_build_Model_collect_interfaces(t *testing.T) {

	// Arrange
	fileVisitor := NewFileVisitor()
	sut := NewSyntaxWalker(fileVisitor)

	source := "" +
		"package main\n\n" +
		"" +
		"type SayHelloFunc func(name string)\n\n" +
		"" +
		"type Greeter interface {\n" +
		"	SayHello(name string)\n" +
		"}"

	fileAccessor := AstFromSource([]byte(source))

	// Act
	err := sut.WalkSyntaxTree(fileAccessor)

	// Assert
	assert.NoError(t, err)

	model, modelErr := fileVisitor.Model()
	assert.NoError(t, modelErr)
	assert.NotNil(t, model)

	assert.Equal(t, "Greeter", model.Interfaces[0].Name)
}

func Test_FileWalker_WalkSyntaxTree_build_Model_collect_func_types(t *testing.T) {

	// Arrange
	fileVisitor := NewFileVisitor()
	sut := NewSyntaxWalker(fileVisitor)

	source := "" +
		"package main\n\n" +
		"" +
		"type SayHelloFunc func(name string)\n\n"

	fileAccessor := AstFromSource([]byte(source))

	// Act
	err := sut.WalkSyntaxTree(fileAccessor)

	// Assert
	assert.NoError(t, err)

	model, modelErr := fileVisitor.Model()
	assert.NoError(t, modelErr)
	assert.NotNil(t, model)

	assert.Equal(t, "SayHelloFunc", model.FuncTypes[0].Name)
}
