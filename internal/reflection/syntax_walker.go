package reflection

import "go/ast"

type walker struct {
	visitor AstVisitor
}

type AstWalker interface {
	WalkSyntaxTree(file *ast.File) error
}

var _ AstWalker = (*walker)(nil)

// NewSyntaxWalker Creates a new AstWalker object.
func NewSyntaxWalker(visitor AstVisitor) AstWalker {
	return &walker{
		visitor: visitor,
	}
}

func (w *walker) WalkSyntaxTree(file *ast.File) error {
	ast.Inspect(file, func(n ast.Node) bool {
		return w.visitor.VisitNode(n)
	})

	return nil
}
