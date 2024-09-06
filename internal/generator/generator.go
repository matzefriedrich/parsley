package generator

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type codeFileGenerator struct {
	options CodeFileGeneratorOptions
}

type CodeFileGeneratorOptions struct {
	TemplateLoader         TemplateLoader
	ConfigureModelCallback ModelConfigurationFunc
	kind                   string
}

type CodeFileGeneratorOptionsFunc func(config *CodeFileGeneratorOptions)

func NewCodeFileGenerator(kind string, config ...CodeFileGeneratorOptionsFunc) (*codeFileGenerator, error) {
	options := CodeFileGeneratorOptions{
		kind: kind,
	}
	for _, f := range config {
		f(&options)
	}
	if options.TemplateLoader == nil {
		return nil, fmt.Errorf("template loader is not set")
	}
	return &codeFileGenerator{options: options}, nil
}

func (g *codeFileGenerator) GenerateCode() {

	goFilePath, err := GetGoFilePath()
	if err != nil {
		fmt.Println(err)
		return
	}

	gen := NewGenericCodeGenerator(g.options.TemplateLoader)

	builder, err := NewTemplateModelBuilder(goFilePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	model, err := builder.Build()
	if err != nil {
		fmt.Println(err)
		return
	}

	if g.options.ConfigureModelCallback != nil {
		g.options.ConfigureModelCallback(model)
	}

	goFileName := path.Base(goFilePath)
	goFileNameWithoutExtension := strings.TrimSuffix(goFileName, filepath.Ext(goFileName))
	goFileDirectory := path.Dir(goFilePath)

	targetFilePath := path.Join(goFileDirectory, fmt.Sprintf("%s.%s.g.go", goFileNameWithoutExtension, g.options.kind))
	f, _ := os.OpenFile(targetFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	defer f.Close()

	gen.Generate(g.options.kind, model, f)

}
