package generator

import (
	"github.com/matzefriedrich/parsley/internal/generator"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_GoFileAccessor_reads_filename_from_GOFILE_environment_variable_returns_error_if_file_not_found(t *testing.T) {

	// Arrange
	const expectedGoFile = "gofile"
	const key = "GOFILE"
	_ = os.Setenv(key, expectedGoFile)

	sut := generator.GoFileAccessor()

	// Act
	source, err := sut()

	// Assert
	assert.ErrorIs(t, err, generator.ErrFailedToObtainGeneratorSourceFile)
	assert.Nil(t, source)
}
