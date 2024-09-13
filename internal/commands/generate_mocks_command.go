package commands

import (
	"fmt"
	"github.com/matzefriedrich/cobra-extensions/pkg"
	"github.com/matzefriedrich/cobra-extensions/pkg/abstractions"
	"github.com/matzefriedrich/parsley/internal/generator"
	"github.com/matzefriedrich/parsley/internal/reflection"
	"github.com/matzefriedrich/parsley/internal/templates"
	"github.com/spf13/cobra"
)

type mocksGeneratorCommand struct {
	use abstractions.CommandName `flag:"mocks" short:"Generate configurable mocks for interface types."`
}

func (m *mocksGeneratorCommand) Execute() {

	templateLoader := func(_ string) (string, error) {
		return templates.MockTemplate, nil
	}

	kind := "mocks"
	gen, _ := generator.NewCodeFileGenerator(kind, func(config *generator.CodeFileGeneratorOptions) {
		config.TemplateLoader = templateLoader
		config.ConfigureModelCallback = func(m *reflection.Model) {
			m.AddImport("github.com/matzefriedrich/parsley/pkg/features")
		}
	})

	err := gen.GenerateCode()
	if err != nil {
		fmt.Println(err)
	}
}

var _ pkg.TypedCommand = (*mockGeneratorCommand)(nil)

func NewGenerateMocksCommand() *cobra.Command {
	command := &mocksGeneratorCommand{}
	return pkg.CreateTypedCommand(command)
}
