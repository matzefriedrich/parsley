package reflection

import (
	"os"
	"testing"

	"github.com/matzefriedrich/parsley/internal/reflection"
	"github.com/stretchr/testify/assert"
)

func Test_AstFromFile_returns_accessor_for_file(t *testing.T) {

	// Arrange
	content := "package main\n\ntype Foo interface {}\n"
	tmpFile, err := os.CreateTemp("", "test*.go")
	assert.NoError(t, err)
	defer func(name string) {
		_ = os.Remove(name)
	}(tmpFile.Name())

	_, err = tmpFile.WriteString(content)
	assert.NoError(t, err)
	_ = tmpFile.Close()

	// Act
	accessor := reflection.AstFromFile(tmpFile.Name())
	file, err := accessor()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, file)
	assert.Equal(t, "main", file.File.Name.Name)
}

func Test_AstFromFile_returns_error_if_file_does_not_exist(t *testing.T) {

	// Arrange
	const path = "non-existent-file.go"

	// Act
	accessor := reflection.AstFromFile(path)
	_, err := accessor()

	// Assert
	assert.Error(t, err)
}
