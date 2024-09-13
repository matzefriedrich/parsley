package reflection

import "go/ast"

type fileCollector struct {
	packageName string
	imports     []string
	interfaces  []Interface
	funcTypes   []FuncType
	comments    []Comment
}

func (t *fileCollector) Model() (*Model, error) {
	return &Model{
		PackageName: t.packageName,
		Imports:     t.imports,
		Interfaces:  t.interfaces,
		FuncTypes:   t.funcTypes,
		Comments:    t.comments,
	}, nil
}

var _ AstFileVisitor = (*fileCollector)(nil)

func NewFileVisitor() AstFileVisitor {
	return &fileCollector{
		imports:    make([]string, 0),
		interfaces: make([]Interface, 0),
		funcTypes:  make([]FuncType, 0),
		comments:   make([]Comment, 0),
	}
}

func (t *fileCollector) VisitNode(n ast.Node) bool {
	switch n.(type) {
	case *ast.CommentGroup:
		commentsGroup, _ := n.(*ast.CommentGroup)
		t.walkComments(commentsGroup)
	case *ast.File:
		f := n.(*ast.File)
		t.VisitFile(f)
	case *ast.ImportSpec:
		spec, _ := n.(*ast.ImportSpec)
		t.VisitImport(spec)
	case *ast.TypeSpec:
		spec, _ := n.(*ast.TypeSpec)
		t.walkTypeSpecNode(spec)
	}
	return true
}

func (t *fileCollector) VisitComment(comment *ast.Comment) {
	model := Comment{
		SymbolInfo: SymbolInfo{
			Pos: comment.Pos(),
			End: comment.End(),
		},
		Text: comment.Text,
	}
	t.comments = append(t.comments, model)
}

func (t *fileCollector) VisitFile(file *ast.File) {
	name := file.Name
	t.packageName = name.Name
}

func (t *fileCollector) VisitImport(importSpec *ast.ImportSpec) {
	t.imports = append(t.imports, importSpec.Path.Value)
}

func (t *fileCollector) VisitInterfaceType(name string, interfaceType *ast.InterfaceType) {
	model := InterfaceWithName(name, SymbolInfo{Pos: interfaceType.Pos(), End: interfaceType.End()})
	methodsCollector := newInterfaceMethodsCollector(&model)
	walker := NewTypeWalker(methodsCollector)
	walker.WalkInterface(interfaceType)
	t.interfaces = append(t.interfaces, model)
}

func (t *fileCollector) VisitFuncType(name string, funcType *ast.FuncType) {
	parameters := CollectParametersFor(funcType)
	results := CollectResultFieldsFor(funcType)
	model := FuncType{
		SymbolInfo: SymbolInfo{
			Pos: funcType.Pos(),
			End: funcType.End(),
		},
		Name:       name,
		Parameters: parameters,
		Results:    results,
	}
	t.funcTypes = append(t.funcTypes, model)
}

func (t *fileCollector) VisitStructType(_ string, _ *ast.StructType) {

}

func (t *fileCollector) walkTypeSpecNode(spec *ast.TypeSpec) {
	typeName := spec.Name.Name
	switch spec.Type.(type) {
	case *ast.InterfaceType:
		interfaceType, _ := spec.Type.(*ast.InterfaceType)
		t.VisitInterfaceType(typeName, interfaceType)
	case *ast.FuncType:
		funcType, _ := spec.Type.(*ast.FuncType)
		t.VisitFuncType(typeName, funcType)
	case *ast.StructType:
		structType, _ := spec.Type.(*ast.StructType)
		t.VisitStructType(typeName, structType)
	}
}

func (t *fileCollector) walkComments(commentGroup *ast.CommentGroup) {
	for _, comment := range commentGroup.List {
		t.VisitComment(comment)
	}
}
