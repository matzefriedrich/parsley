package commands

import (
	"fmt"
	"github.com/matzefriedrich/parsley/internal/templates"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/matzefriedrich/cobra-extensions/pkg"
	"github.com/matzefriedrich/cobra-extensions/pkg/abstractions"
	"github.com/matzefriedrich/parsley/internal/generator"
	"github.com/spf13/cobra"
)

type generateProxyCommand struct {
	use abstractions.CommandName `flag:"proxy" short:"GenerateProjectFiles generic proxy types for method call interception."`
}

func (g *generateProxyCommand) Execute() {

	goFilePath, err := generator.GetGoFilePath()
	if err != nil {
		fmt.Println(err)
		return
	}

	gen := generator.NewGenericCodeGenerator(func(_ string) (string, error) {
		return templates.ProxyTemplate, nil
	})

	builder, err := generator.NewTemplateModelBuilder(goFilePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	model, err := builder.Build()
	if err != nil {
		fmt.Println(err)
		return
	}

	model.AddImport("github.com/matzefriedrich/parsley/pkg/features")

	goFileName := path.Base(goFilePath)
	goFileNameWithoutExtension := strings.TrimSuffix(goFileName, filepath.Ext(goFileName))
	goFileDirectory := path.Dir(goFilePath)

	targetFilePath := path.Join(goFileDirectory, fmt.Sprintf("%s.proxy.g.go", goFileNameWithoutExtension))
	f, _ := os.OpenFile(targetFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	defer f.Close()

	gen.Generate("proxy", model, f)
}

var _ pkg.TypedCommand = &generateProxyCommand{}

func NewGenerateProxyCommand() *cobra.Command {
	command := &generateProxyCommand{}
	return pkg.CreateTypedCommand(command)
}
