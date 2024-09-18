package mocks

import (
	"github.com/matzefriedrich/parsley/internal/generator"
	"io"
)

type GenericCodeGeneratorMock struct {
	AddTemplateFuncFunc func(functions ...generator.TemplateFunction) error
	GenerateFunc        func(templateName string, model any, writer io.Writer) error
}

func NewGenericCodeGeneratorMock() *GenericCodeGeneratorMock {
	return &GenericCodeGeneratorMock{}
}

func (g *GenericCodeGeneratorMock) AddTemplateFunc(functions ...generator.TemplateFunction) error {
	return g.AddTemplateFuncFunc(functions...)
}

func (g *GenericCodeGeneratorMock) Generate(templateName string, model any, writer io.Writer) error {
	return g.GenerateFunc(templateName, model, writer)
}

var _ generator.GenericCodeGenerator = (*GenericCodeGeneratorMock)(nil)
