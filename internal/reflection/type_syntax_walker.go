package reflection

import "go/ast"

type typeWalker struct {
	visitor AstTypeVisitor
}

type AstTypeSpecWalker interface {
	WalkInterface(interfaceType *ast.InterfaceType)
	WalkFunc(name string, funcSpec *ast.FuncType)
}

var _ AstTypeSpecWalker = (*typeWalker)(nil)

// NewTypeWalker Creates a new AstTypeSpecWalker object.
func NewTypeWalker(typeVisitor AstTypeVisitor) AstTypeSpecWalker {
	return &typeWalker{
		visitor: typeVisitor,
	}
}

func (t *typeWalker) WalkInterface(interfaceType *ast.InterfaceType) {
	for _, method := range interfaceType.Methods.List {
		methodNames := method.Names
		if methodNames == nil || len(methodNames) == 0 {
			continue
		}
		name := methodNames[0].Name
		if funcType, ok := method.Type.(*ast.FuncType); ok {
			t.WalkFunc(name, funcType)
		}
	}
}

func (t *typeWalker) WalkFunc(name string, funcSpec *ast.FuncType) {
	t.visitor.VisitMethod(name, funcSpec)
}
