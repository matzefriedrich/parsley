package reflection

import "go/ast"

type fileWalker struct {
	visitor AstFileVisitor
}

type AstFileWalker interface {
	WalkSyntaxTree(accessor AstFileAccessor) error
}

var _ AstFileWalker = (*fileWalker)(nil)

// NewFileWalker Creates a new AstFileWalker object.
func NewFileWalker(visitor AstFileVisitor) AstFileWalker {
	return &fileWalker{
		visitor: visitor,
	}
}

func (g *fileWalker) WalkSyntaxTree(accessor AstFileAccessor) error {

	file, err := accessor()
	if err != nil {
		return err
	}

	ast.Inspect(file, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.File:
			f := n.(*ast.File)
			g.visitor.VisitFile(f)
		case *ast.ImportSpec:
			spec, _ := n.(*ast.ImportSpec)
			g.visitor.VisitImport(spec)
		case *ast.TypeSpec:
			spec, _ := n.(*ast.TypeSpec)
			g.walkTypeSpecNode(spec)
		}
		return true
	})

	return nil
}

func (g *fileWalker) walkTypeSpecNode(spec *ast.TypeSpec) {
	typeName := spec.Name.Name
	switch spec.Type.(type) {
	case *ast.InterfaceType:
		interfaceType, _ := spec.Type.(*ast.InterfaceType)
		g.visitor.VisitInterfaceType(typeName, interfaceType)
	case *ast.FuncType:
		funcType, _ := spec.Type.(*ast.FuncType)
		g.visitor.VisitFuncType(typeName, funcType)
	case *ast.StructType:
		structType, _ := spec.Type.(*ast.StructType)
		g.visitor.VisitStructType(typeName, structType)
	}
}
