package commands

import (
	"fmt"
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
	p, err := generator.OpenProject(projectFolderPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	const minVersion = "v0.8.1"
	const packageName = "github.com/matzefriedrich/parsley"
	dependencyErr := p.AddDependency(packageName, minVersion)
	if dependencyErr != nil {
		fmt.Println(err)
		return
	}

	gen := generator.NewBootstrapGenerator(projectFolderPath)
	gen.GenerateProjectFiles()
}

var _ pkg.TypedCommand = &initCommand{}

func NewInitCommand() *cobra.Command {
	command := &initCommand{}
	return pkg.CreateTypedCommand(command)
}
