package reflection

import (
	"github.com/matzefriedrich/parsley/internal/reflection"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Model_AddImport_adds_import(t *testing.T) {
	// Arrange
	sut := &reflection.Model{Imports: make([]string, 0)}
	hasNoImports := len(sut.Imports) == 0

	// Act
	sut.AddImport("fmt")
	actual := sut.Imports[0]

	// Assert
	assert.True(t, hasNoImports)
	assert.Equal(t, 1, len(sut.Imports))
	assert.Equal(t, "fmt", actual)
}

func Test_Model_AddImport_adds_quoted_import(t *testing.T) {
	// Arrange
	sut := &reflection.Model{Imports: make([]string, 0)}
	hasNoImports := len(sut.Imports) == 0

	// Act
	sut.AddImport("\"fmt\"")
	actual := sut.Imports[0]

	// Assert
	assert.True(t, hasNoImports)
	assert.Equal(t, 1, len(sut.Imports))
	assert.Equal(t, "fmt", actual)
}

func Test_Parameter_MatchesType_returns_true_if_parameter_type_matches_expected_typename(t *testing.T) {
	// Arrange
	sut := &reflection.Parameter{Name: "p", TypeName: "string"}
	// Act
	actual := sut.MatchesType("string")
	// Assert
	assert.True(t, actual)
}

func Test_Parameter_MatchesType_returns_false_if_type_does_not_match_expected_typename(t *testing.T) {
	// Arrange
	sut := &reflection.Parameter{Name: "p", TypeName: "string"}
	// Act
	actual := sut.MatchesType("bool")
	// Assert
	assert.False(t, actual)
}
