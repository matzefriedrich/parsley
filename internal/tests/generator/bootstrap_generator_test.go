package generator

import (
	"github.com/matzefriedrich/parsley/internal/generator"
	"github.com/matzefriedrich/parsley/internal/tests/mocks"
	"io"
	"testing"
)

func Test_BootstrapGenerator_ScaffoldProjectFiles_generates_expected_project_files(t *testing.T) {

	// Arrange
	expectedProjectFiles := []string{"application.go", "main.go", "greeter.go"}
	memoryFiles := make(map[string]mocks.MemoryFile)
	writerFuncFactory := func(targetFilename string) (io.WriteCloser, error) {
		f, found := memoryFiles[targetFilename]
		if found {
			return f, nil
		}
		f = mocks.NewMemoryFile()
		memoryFiles[targetFilename] = f
		return f, nil
	}

	sut := generator.NewBootstrapGenerator(writerFuncFactory)

	// Act
	sut.ScaffoldProjectFiles()

	// Assert
	for _, expectedProjectFile := range expectedProjectFiles {
		file, found := memoryFiles[expectedProjectFile]
		if !found {
			t.Errorf("Expected project file %s not found", expectedProjectFile)
		}
		content := file.String()
		if len(content) == 0 {
			t.Errorf("Expected project file %s to have content", expectedProjectFile)
		}
		t.Logf("Project file %s has content:\n%s", expectedProjectFile, content)
	}
}
