package commands

import (
	"github.com/matzefriedrich/cobra-extensions/pkg"
	"github.com/matzefriedrich/cobra-extensions/pkg/abstractions"
	"github.com/matzefriedrich/parsley/internal/generator"
	"github.com/spf13/cobra"
	"os"
)

type initCommand struct {
	use abstractions.CommandName `flag:"init" short:"Add Parsley to an application"`
}

func (g *initCommand) Execute() {
	projectFolderPath, _ := os.Getwd()
	gen := generator.NewBootstrapGenerator(projectFolderPath)
	gen.GenerateProjectFiles()
}

var _ pkg.TypedCommand = &initCommand{}

func NewInitCommand() *cobra.Command {
	command := &initCommand{}
	return pkg.CreateTypedCommand(command)
}
