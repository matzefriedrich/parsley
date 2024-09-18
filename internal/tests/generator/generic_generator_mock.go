package generator

import (
	"github.com/matzefriedrich/parsley/internal/generator"
	"io"
)

type genericCodeGeneratorMock struct {
	AddTemplateFuncFunc func(functions ...generator.TemplateFunction) error
	GenerateFunc        func(templateName string, model any, writer io.Writer) error
}

func (g *genericCodeGeneratorMock) AddTemplateFunc(functions ...generator.TemplateFunction) error {
	return g.AddTemplateFuncFunc(functions...)
}

func (g *genericCodeGeneratorMock) Generate(templateName string, model any, writer io.Writer) error {
	return g.GenerateFunc(templateName, model, writer)
}

var _ generator.GenericCodeGenerator = (*genericCodeGeneratorMock)(nil)
