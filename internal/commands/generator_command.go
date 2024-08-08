package commands

import (
	"github.com/matzefriedrich/cobra-extensions/pkg"
	"github.com/matzefriedrich/cobra-extensions/pkg/abstractions"
	"github.com/spf13/cobra"
)

type generatorCommand struct {
	use abstractions.CommandName `flag:"generate"`
}

func (g *generatorCommand) Execute() {

}

var _ pkg.TypedCommand = &generatorCommand{}

func NewGenerateGroupCommand() *cobra.Command {
	command := &generatorCommand{}
	return pkg.CreateTypedCommand(command)
}
