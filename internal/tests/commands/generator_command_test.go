package commands

import (
	"github.com/matzefriedrich/parsley/internal/commands"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_NewGenerateGroupCommand_Execute(t *testing.T) {

	// Arrange
	sut := commands.NewGenerateGroupCommand()

	// Act
	err := sut.Execute()

	// Assert
	assert.NoError(t, err)
}
