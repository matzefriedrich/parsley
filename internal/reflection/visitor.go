package reflection

import (
	"go/ast"
)

type AstFileVisitor interface {
	VisitFile(file *ast.File)
	VisitImport(importSpec *ast.ImportSpec)
	VisitInterfaceType(name string, interfaceType *ast.InterfaceType)
	VisitFuncType(name string, funcType *ast.FuncType)
	VisitStructType(name string, structType *ast.StructType)
	Model() (*Model, error)
}

type AstTypeVisitor interface {
	VisitMethod(name string, method *ast.FuncType)
}
