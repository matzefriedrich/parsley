package commands

import (
	"github.com/matzefriedrich/parsley/internal/commands"
	"github.com/matzefriedrich/parsley/internal/reflection"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func Test_GenerateMocksCommand_Execute(t *testing.T) {

	// Arrange
	source := []byte("package main\n" + "\n" +
		"//parsley:mock" + "\n" +
		"type Greeter interface {\n" +
		"	SayHello(name string)" + "\n" +
		"}")

	buffer := newMemoryFile()
	outputWriterFactory := func(kind string, source *reflection.AstFileSource) (io.WriteCloser, error) {
		return buffer, nil
	}

	fileAccessor := reflection.AstFromSource(source)
	sut := commands.NewGenerateMocksCommand(fileAccessor, outputWriterFactory)

	// Act
	err := sut.Execute()

	actual := buffer.String()

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, actual)
}
