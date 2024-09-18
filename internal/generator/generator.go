package generator

import (
	"fmt"
	"github.com/matzefriedrich/parsley/internal/reflection"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type codeFileGenerator struct {
	options      CodeFileGeneratorOptions
	fileAccessor reflection.AstFileAccessor
}

type CodeFileGeneratorOptions struct {
	TemplateLoader         TemplateLoader
	ConfigureModelCallback reflection.ModelConfigurationFunc
	OutputWriterFactory    OutputWriterFactory
	kind                   string
}

type CodeFileGenerator interface {
	GenerateCode() error
}

type CodeFileGeneratorOptionsFunc func(config *CodeFileGeneratorOptions)

// FileOutputWriter Creates an OutputWriterFactory object that can be used create file writers.
func FileOutputWriter() OutputWriterFactory {
	return func(kind string, source *reflection.AstFileSource) (OutputWriter, error) {
		fileName := path.Base(source.Filename)
		fileNameWithoutExtension := strings.TrimSuffix(fileName, filepath.Ext(fileName))
		fileDirectory := path.Dir(source.Filename)

		targetFilePath := path.Join(fileDirectory, fmt.Sprintf("%s.%s.g.go", fileNameWithoutExtension, kind))
		return os.OpenFile(targetFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	}
}

func NewCodeFileGenerator(kind string, fileAccessor reflection.AstFileAccessor, config ...CodeFileGeneratorOptionsFunc) (CodeFileGenerator, error) {
	options := CodeFileGeneratorOptions{
		kind: kind,
	}
	for _, f := range config {
		f(&options)
	}
	if options.TemplateLoader == nil {
		return nil, fmt.Errorf("template loader is not set")
	}
	return &codeFileGenerator{
		fileAccessor: fileAccessor,
		options:      options,
	}, nil
}

func (g *codeFileGenerator) GenerateCode() error {

	gen := NewGenericCodeGenerator(g.options.TemplateLoader)
	err := RegisterTemplateFunctions(gen, RegisterTypeModelFunctions, RegisterNamingFunctions)
	if err != nil {
		return err
	}

	source, err := g.fileAccessor()
	if err != nil {
		return err
	}

	builder := NewTemplateModelBuilder(source.File)

	model, err := builder.Build()
	if err != nil {
		return err
	}

	if g.options.ConfigureModelCallback != nil {
		g.options.ConfigureModelCallback(model)
	}

	f, outputErr := g.options.OutputWriterFactory(g.options.kind, source)
	if outputErr != nil {
		return outputErr
	}

	defer func(f OutputWriter) {
		_ = f.Close()
	}(f)

	err = gen.Generate(g.options.kind, model, f)
	if err != nil {
		return err
	}

	return nil
}
