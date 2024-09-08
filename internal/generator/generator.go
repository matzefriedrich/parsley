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

type CodeFileGenerator interface {
	GenerateCode() error
}

type CodeFileGeneratorOptionsFunc func(config *CodeFileGeneratorOptions)

func NewCodeFileGenerator(kind string, config ...CodeFileGeneratorOptionsFunc) (CodeFileGenerator, error) {
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

func (g *codeFileGenerator) GenerateCode() error {

	goFilePath, err := GetGoFilePath()
	if err != nil {
		return err
	}

	gen := NewGenericCodeGenerator(g.options.TemplateLoader)
	err = RegisterTemplateFunctions(gen, RegisterTypeModelFunctions, RegisterNamingFunctions)
	if err != nil {
		return err
	}

	builder, err := NewTemplateModelBuilder(goFilePath)
	if err != nil {
		return err
	}

	model, err := builder.Build()
	if err != nil {
		return err
	}

	if g.options.ConfigureModelCallback != nil {
		g.options.ConfigureModelCallback(model)
	}

	goFileName := path.Base(goFilePath)
	goFileNameWithoutExtension := strings.TrimSuffix(goFileName, filepath.Ext(goFileName))
	goFileDirectory := path.Dir(goFilePath)

	targetFilePath := path.Join(goFileDirectory, fmt.Sprintf("%s.%s.g.go", goFileNameWithoutExtension, g.options.kind))
	f, _ := os.OpenFile(targetFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	err = gen.Generate(g.options.kind, model, f)
	if err != nil {
		return err
	}

	return nil
}
