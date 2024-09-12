package generator

import (
	"github.com/matzefriedrich/parsley/internal/reflection"
)

type TemplateModelBuilder struct {
	accessor reflection.AstFileAccessor
}

func NewTemplateModelBuilder(accessor reflection.AstFileAccessor) *TemplateModelBuilder {
	return &TemplateModelBuilder{
		accessor: accessor,
	}
}

func (b *TemplateModelBuilder) Build() (*reflection.Model, error) {

	fileVisitor := reflection.NewFileVisitor()
	walker := reflection.NewFileWalker(fileVisitor)
	err := walker.WalkSyntaxTree(b.accessor)
	if err != nil {
		return nil, err
	}

	return fileVisitor.Model()
}
