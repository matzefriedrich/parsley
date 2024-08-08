package generator

import (
	"github.com/pkg/errors"
	"io"
	"os"
	"text/template"
)

type generator struct {
	templateLoader TemplateLoader
}

type GenericCodeGenerator interface {
	Generate(m *Model, writer io.Writer) error
}

type TemplateLoader func() (string, error)

func NewGenericCodeGenerator(templateLoader TemplateLoader) GenericCodeGenerator {
	return &generator{
		templateLoader: templateLoader,
	}
}

func (g *generator) Generate(m *Model, writer io.Writer) error {
	tmpl, err := g.templateLoader()
	if err != nil {
		return errors.Wrap(err, ErrorCannotGenerateProxies)
	}

	t := template.Must(template.New("").Parse(tmpl))
	err = t.Execute(writer, m)
	if err != nil {
		return errors.Wrap(err, ErrorCannotExecuteTemplate)
	}

	return nil
}

func (g *generator) LoadTemplateFromFile(templateFile string) (string, error) {

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
