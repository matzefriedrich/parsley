package generator

import (
	"github.com/matzefriedrich/parsley/internal/reflection"
	"go/ast"
)

type TemplateModelBuilder struct {
	file *ast.File
}

func NewTemplateModelBuilder(file *ast.File) *TemplateModelBuilder {
	return &TemplateModelBuilder{
		file: file,
	}
}

func (b *TemplateModelBuilder) Build() (*reflection.Model, error) {

	fileVisitor := reflection.NewFileVisitor()
	walker := reflection.NewSyntaxWalker(fileVisitor)
	err := walker.WalkSyntaxTree(b.file)
	if err != nil {
		return nil, err
	}

	return fileVisitor.Model()
}
