package commands

import (
	"github.com/matzefriedrich/cobra-extensions/pkg/commands"
	"github.com/matzefriedrich/cobra-extensions/pkg/types"
	"github.com/spf13/cobra"
)

type generatorCommand struct {
	use types.CommandName `flag:"generate" short:"Generate boilerplate code for advanced DI features."`
}

func (g *generatorCommand) Execute() {

}

var _ types.TypedCommand = &generatorCommand{}

func NewGenerateGroupCommand() *cobra.Command {
	command := &generatorCommand{}
	return commands.CreateTypedCommand(command)
}
