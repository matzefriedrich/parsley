package commands

import (
	"github.com/matzefriedrich/parsley/internal/commands"
	"github.com/matzefriedrich/parsley/internal/generator"
	"github.com/matzefriedrich/parsley/internal/tests/mocks"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func Test_InitCommand_Execute_adds_project_reference_and_scaffolds_files(t *testing.T) {

	// Arrange
	expectedProjectFiles := []string{"application.go", "main.go", "greeter.go"}
	files := make(map[string]mocks.MemoryFile)

	writerFactory := func(projectFolder string) (generator.ScaffoldingFileWriterFunc, error) {
		return func(targetFilename string) (io.WriteCloser, error) {
			f, found := files[targetFilename]
			if found {
				return f, nil
			}
			f = mocks.NewMemoryFile()
			files[targetFilename] = f
			return f, nil
		}, nil
	}

	projectInstance := &memoryGoProject{
		packages: make(map[string]string),
	}

	sut := commands.NewInitCommand(writerFactory, func(projectFolderPath string) (generator.GoProject, error) {
		return projectInstance, nil
	})

	// Act
	err := sut.Execute()

	// Assert
	assert.NoError(t, err)

	_, hasParsleyReference := projectInstance.packages["github.com/matzefriedrich/parsley"]
	assert.True(t, hasParsleyReference)

	for _, expectedProjectFile := range expectedProjectFiles {
		_, found := files[expectedProjectFile]
		if !found {
			t.Errorf("Expected project file %s not found", expectedProjectFile)
		}
	}
}

type memoryGoProject struct {
	packages map[string]string
}

func (m *memoryGoProject) AddDependency(packageName string, version string) error {
	m.packages[packageName] = version
	return nil
}

var _ generator.GoProject = (*memoryGoProject)(nil)
