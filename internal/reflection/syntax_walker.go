package reflection

import "go/ast"

type walker struct {
	visitor AstVisitor
}

type AstWalker interface {
	WalkSyntaxTree(accessor AstFileAccessor) error
}

var _ AstWalker = (*walker)(nil)

// NewSyntaxWalker Creates a new AstWalker object.
func NewSyntaxWalker(visitor AstVisitor) AstWalker {
	return &walker{
		visitor: visitor,
	}
}

func (w *walker) WalkSyntaxTree(accessor AstFileAccessor) error {
	file, err := accessor()
	if err != nil {
		return err
	}

	ast.Inspect(file, func(n ast.Node) bool {
		return w.visitor.VisitNode(n)
	})

	return nil
}
