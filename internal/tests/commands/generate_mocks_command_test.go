package commands

import (
	"github.com/matzefriedrich/parsley/internal/commands"
	"github.com/matzefriedrich/parsley/internal/generator"
	"github.com/matzefriedrich/parsley/internal/reflection"
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
)

func Test_GenerateMocksCommand_Execute(t *testing.T) {

	// Arrange
	source := []byte("package main\n" + "\n" +
		"type Greeter interface {\n" +
		"	SayHello(name string)" + "\n" +
		"}")

	buffer := memoryFileTarget{buffer: strings.Builder{}}
	outputWriterFactory := func(kind string, source *reflection.AstFileSource) (generator.OutputWriter, error) {
		return &buffer, nil
	}

	fileAccessor := reflection.AstFromSource(source)
	sut := commands.NewGenerateMocksCommand(fileAccessor, outputWriterFactory)

	// Act
	err := sut.Execute()

	actual := buffer.buffer.String()

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, actual)
}

type memoryFileTarget struct {
	buffer strings.Builder
}

func (m *memoryFileTarget) Write(p []byte) (n int, err error) {
	return m.buffer.Write(p)
}

func (m *memoryFileTarget) Close() error {
	return nil
}

var _ io.Writer = (*memoryFileTarget)(nil)
var _ io.Closer = (*memoryFileTarget)(nil)
