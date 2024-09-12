package reflection

import "go/ast"

type fileCollector struct {
	packageName string
	imports     []string
	interfaces  []Interface
	funcTypes   []FuncType
}

func (t *fileCollector) Model() (*Model, error) {
	return &Model{
		PackageName: t.packageName,
		Imports:     t.imports,
		Interfaces:  t.interfaces,
		FuncTypes:   t.funcTypes,
	}, nil
}

var _ AstFileVisitor = (*fileCollector)(nil)

func NewFileVisitor() AstFileVisitor {
	return &fileCollector{
		imports:    make([]string, 0),
		interfaces: make([]Interface, 0),
		funcTypes:  make([]FuncType, 0),
	}
}

func (t *fileCollector) VisitFile(file *ast.File) {
	name := file.Name
	t.packageName = name.Name
}

func (t *fileCollector) VisitImport(importSpec *ast.ImportSpec) {
	t.imports = append(t.imports, importSpec.Path.Value)
}

func (t *fileCollector) VisitInterfaceType(name string, interfaceType *ast.InterfaceType) {
	model := InterfaceWithName(name)
	methodsCollector := newInterfaceMethodsCollector(&model)
	walker := NewTypeWalker(methodsCollector)
	walker.WalkInterface(interfaceType)
	t.interfaces = append(t.interfaces, model)
}

func (t *fileCollector) VisitFuncType(name string, funcType *ast.FuncType) {
	parameters := CollectParametersFor(funcType)
	results := CollectResultFieldsFor(funcType)
	model := FuncType{
		Name:       name,
		Parameters: parameters,
		Results:    results,
	}
	t.funcTypes = append(t.funcTypes, model)
}

func (t *fileCollector) VisitStructType(_ string, _ *ast.StructType) {

}
