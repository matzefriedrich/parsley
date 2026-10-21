package commands

import (
	"io"
	"testing"

	"github.com/matzefriedrich/parsley/internal/commands"
	"github.com/matzefriedrich/parsley/internal/reflection"
	"github.com/matzefriedrich/parsley/internal/tests/mocks"
	"github.com/stretchr/testify/assert"
)

func Test_GenerateMocksCommand_Execute(t *testing.T) {

	// Arrange
	source := []byte("package main\n" + "\n" +
		"//parsley:mock" + "\n" +
		"type Greeter interface {\n" +
		"	SayHello(name string)" + "\n" +
		"}")

	buffer := mocks.NewMemoryFile()
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

func Test_GenerateMocksCommand_Execute_uses_only_required_imports_and_aliases(t *testing.T) {
	source := []byte("package main\n" +
		"import (\n\tunused \"example.org/unused\"\n\twidgets \"example.org/widgets\"\n)\n" +
		"//parsley:mock\n" +
		"type Service interface {\n\tLoad() widgets.Widget\n}\n")

	buffer := mocks.NewMemoryFile()
	outputWriterFactory := func(_ string, _ *reflection.AstFileSource) (io.WriteCloser, error) {
		return buffer, nil
	}
	command := commands.NewGenerateMocksCommand(reflection.AstFromSource(source), outputWriterFactory)

	err := command.Execute()
	generated := buffer.String()

	assert.NoError(t, err)
	assert.Contains(t, generated, `widgets "example.org/widgets"`)
	assert.Contains(t, generated, "widgets.Widget")
	assert.NotContains(t, generated, "example.org/unused")
}

func Test_GenerateMocksCommand_Execute_omits_imports_when_all_interfaces_filtered(t *testing.T) {
	source := []byte("package main\n" +
		"//parsley:ignore\n" +
		"type Service interface {\n\tLoad()\n}\n")

	buffer := mocks.NewMemoryFile()
	outputWriterFactory := func(_ string, _ *reflection.AstFileSource) (io.WriteCloser, error) {
		return buffer, nil
	}
	command := commands.NewGenerateMocksCommand(reflection.AstFromSource(source), outputWriterFactory)

	err := command.Execute()
	generated := buffer.String()

	assert.NoError(t, err)
	assert.NotContains(t, generated, "import (")
	assert.NotContains(t, generated, "pkg/features")
}
