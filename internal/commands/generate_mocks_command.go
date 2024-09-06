package commands

import (
	"github.com/matzefriedrich/cobra-extensions/pkg"
	"github.com/matzefriedrich/cobra-extensions/pkg/abstractions"
	"github.com/matzefriedrich/parsley/internal/generator"
	"github.com/matzefriedrich/parsley/internal/templates"
	"github.com/spf13/cobra"
)

type mockGeneratorCommand struct {
	use abstractions.CommandName `flag:"mocks" short:"Generate configurable mocks for interface types."`
}

func (m *mockGeneratorCommand) Execute() {

	templateLoader := func(_ string) (string, error) {
		return templates.MockTemplate, nil
	}

	kind := "mocks"
	gen, _ := generator.NewCodeFileGenerator(kind, func(config *generator.CodeFileGeneratorOptions) {
		config.TemplateLoader = templateLoader
	})

	gen.GenerateCode()
}

var _ pkg.TypedCommand = (*mockGeneratorCommand)(nil)

func NewGenerateMocksCommand() *cobra.Command {
	command := &mockGeneratorCommand{}
	return pkg.CreateTypedCommand(command)
}
