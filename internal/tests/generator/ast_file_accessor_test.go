package generator

import (
	"github.com/matzefriedrich/parsley/internal/reflection"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_AstFileAccessor_create_accessor_from_source(t *testing.T) {

	// Arrange
	source := []byte("package types\n" +
		"\n" +
		"type Service0 interface {\n" +
		"	Method0()\n" +
		"}\n" +
		"\n" +
		"type Service1 interface {\n" +
		"	Method1() string\n" +
		"	Method3() (string, error)\n" +
		"}\n")

	accessor := reflection.AstFromSource(source)

	// Act
	actual, err := accessor()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, actual)
}
