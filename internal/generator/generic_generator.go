package generator

import (
	"bytes"
	"github.com/pkg/errors"
	"go/format"
	"io"
	"os"
	"reflect"
	"text/template"
)

type genericGenerator struct {
	templateLoader TemplateLoader
	funcMap        template.FuncMap
}

type GenericCodeGenerator interface {
	AddTemplateFunc(functions ...TemplateFunction) error
	Generate(templateName string, model any, writer io.Writer) error
}

var _ GenericCodeGenerator = (*genericGenerator)(nil)

func NewGenericCodeGenerator(templateLoader TemplateLoader) GenericCodeGenerator {
	g := &genericGenerator{
		templateLoader: templateLoader,
		funcMap:        template.FuncMap{},
	}
	return g
}

func (g *genericGenerator) AddTemplateFunc(functions ...TemplateFunction) error {

	addFunc := func(name string, f any) error {

		if len(name) == 0 {
			return errors.New("function name cannot be empty")
		}
		reflected := reflect.ValueOf(f)
		if reflected.Kind() != reflect.Func {
			return errors.New("the given value is not a function")
		}

		g.funcMap[name] = f
		return nil
	}

	for _, function := range functions {
		if err := addFunc(function.Name, function.Function); err != nil {
			return err
		}
	}

	return nil
}

func (g *genericGenerator) Generate(templateName string, templateModel any, writer io.Writer) error {

	tmpl, err := g.templateLoader(templateName)
	if err != nil {
		return newGeneratorError(ErrorCannotGenerateProxies, WithCause(err))
	}

	var generatedCode bytes.Buffer

	t := template.Must(template.New("").Funcs(g.funcMap).Parse(tmpl))
	err = t.Execute(&generatedCode, templateModel)
	if err != nil {
		return newGeneratorError(ErrorCannotExecuteTemplate, WithCause(err))
	}

	formattedCode, formatErr := format.Source(generatedCode.Bytes())
	if formatErr != nil {
		return newGeneratorError(ErrorCannotFormatGeneratedCode, WithCause(formatErr))
	}

	_, writerErr := writer.Write(formattedCode)
	if writerErr != nil {
		return newGeneratorError(ErrorFailedToWriteGeneratedCode, WithCause(writerErr))
	}

	return nil
}

func (g *genericGenerator) LoadTemplateFromFile(templateFile string) (string, error) {

	if _, err := os.Stat(templateFile); errors.Is(err, os.ErrNotExist) {
		return "", newGeneratorError(ErrorTemplateFileNotFound, WithCause(err))
	}

	f, err := os.OpenFile(templateFile, os.O_RDONLY, 400)
	defer func(file *os.File) {
		_ = file.Close()
	}(f)

	if err != nil {
		return "", newGeneratorError(ErrorFailedToOpenTemplateFile, WithCause(err))
	}

	bytes, _ := io.ReadAll(f)
	return string(bytes), nil
}

func RegisterTemplateFunctions(g GenericCodeGenerator, functions ...func(generator GenericCodeGenerator) error) error {
	for _, function := range functions {
		err := function(g)
		if err != nil {
			return err
		}
	}
	return nil
}
