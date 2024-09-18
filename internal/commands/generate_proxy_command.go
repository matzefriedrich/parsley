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

type generateProxyCommand struct {
	use abstractions.CommandName `flag:"proxy" short:"Generate generic proxy types for method call interception."`
}

func (g *generateProxyCommand) Execute() {

	templateLoader := func(_ string) (string, error) {
		return templates.ProxyTemplate, nil
	}

	kind := "proxy"
	accessor := generator.GoFileAccessor()
	gen, _ := generator.NewCodeFileGenerator(kind, accessor, func(config *generator.CodeFileGeneratorOptions) {
		config.TemplateLoader = templateLoader
		config.OutputWriterFactory = generator.FileOutputWriter()
		config.ConfigureModelCallback = func(m *reflection.Model) {
			m.AddImport("github.com/matzefriedrich/parsley/pkg/features")
		}
	})

	err := gen.GenerateCode()
	if err != nil {
		fmt.Println(err)
	}
}

var _ pkg.TypedCommand = &generateProxyCommand{}

func NewGenerateProxyCommand() *cobra.Command {
	command := &generateProxyCommand{}
	return pkg.CreateTypedCommand(command)
}
