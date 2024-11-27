package commands

import (
	"fmt"
	"github.com/matzefriedrich/cobra-extensions/pkg/commands"
	"github.com/matzefriedrich/cobra-extensions/pkg/types"
	"github.com/matzefriedrich/parsley/internal/generator"
	"github.com/matzefriedrich/parsley/internal/reflection"
	"github.com/matzefriedrich/parsley/internal/templates"
	"github.com/spf13/cobra"
)

type generateProxyCommand struct {
	use                 types.CommandName `flag:"proxy" short:"Generate generic proxy types for method call interception." long:"Generates generic proxy types designed for method call interception on Go interfaces. These proxies act as intermediaries, allowing you to inject custom behavior—such as logging, validation, or transformation—before or after method execution."`
	fileAccessor        reflection.AstFileAccessor
	outputWriterFactory generator.OutputWriterFactory
}

// Execute generates the code for a proxy.
func (g *generateProxyCommand) Execute() {

	templateLoader := func(_ string) (string, error) {
		return templates.ProxyTemplate, nil
	}

	kind := "proxy"
	gen, _ := generator.NewCodeFileGenerator(kind, g.fileAccessor, func(config *generator.CodeFileGeneratorOptions) {
		config.TemplateLoader = templateLoader
		config.OutputWriterFactory = g.outputWriterFactory
		config.ConfigureModelCallback = func(m *reflection.Model) {
			m.AddImport("github.com/matzefriedrich/parsley/pkg/features")
		}
	})

	err := gen.GenerateCode()
	if err != nil {
		fmt.Println(err)
	}
}

var _ types.TypedCommand = &generateProxyCommand{}

// NewGenerateProxyCommand creates a new cobra.Command for generating proxy code, enabling method call interception for interfaces.
func NewGenerateProxyCommand(fileAccessor reflection.AstFileAccessor, outputWriterFactory generator.OutputWriterFactory) *cobra.Command {
	command := &generateProxyCommand{
		fileAccessor:        fileAccessor,
		outputWriterFactory: outputWriterFactory,
	}
	return commands.CreateTypedCommand(command)
}
