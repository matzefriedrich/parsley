package reflection

import "go/ast"

type fileVisitor struct {
	idSequence  uint64
	packageName string
	imports     []string
	interfaces  []Interface
	funcTypes   []FuncType
	comments    []Comment
}

func (t *fileVisitor) Model() (*Model, error) {
	return &Model{
		PackageName: t.packageName,
		Imports:     t.imports,
		Interfaces:  t.interfaces,
		FuncTypes:   t.funcTypes,
		Comments:    t.comments,
	}, nil
}

var _ AstFileVisitor = (*fileVisitor)(nil)

func NewFileVisitor() AstFileVisitor {
	idSeed := 1
	return &fileVisitor{
		idSequence: uint64(idSeed),
		imports:    make([]string, 0),
		interfaces: make([]Interface, 0),
		funcTypes:  make([]FuncType, 0),
		comments:   make([]Comment, 0),
	}
}

func (t *fileVisitor) VisitNode(n ast.Node) bool {
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

func (t *fileVisitor) VisitComment(comment *ast.Comment) {
	model := Comment{
		SymbolInfo: SymbolInfo{
			Pos: comment.Pos(),
			End: comment.End(),
		},
		Text: comment.Text,
	}
	t.comments = append(t.comments, model)
}

func (t *fileVisitor) VisitFile(file *ast.File) {
	name := file.Name
	t.packageName = name.Name
}

func (t *fileVisitor) VisitImport(importSpec *ast.ImportSpec) {
	t.imports = append(t.imports, importSpec.Path.Value)
}

func (t *fileVisitor) VisitInterfaceType(name string, interfaceType *ast.InterfaceType) {
	id := t.newSymbolId()
	model := InterfaceWithName(name, SymbolInfo{Id: id, Pos: interfaceType.Pos(), End: interfaceType.End()})
	methodsCollector := newInterfaceMethodsCollector(&model)
	walker := NewTypeWalker(methodsCollector)
	walker.WalkInterface(interfaceType)
	t.interfaces = append(t.interfaces, model)
}

func (t *fileVisitor) VisitFuncType(name string, funcType *ast.FuncType) {
	parameters := CollectParametersFor(funcType)
	results := CollectResultFieldsFor(funcType)
	model := FuncType{
		SymbolInfo: SymbolInfo{
			Id:  t.newSymbolId(),
			Pos: funcType.Pos(),
			End: funcType.End(),
		},
		Name:       name,
		Parameters: parameters,
		Results:    results,
	}
	t.funcTypes = append(t.funcTypes, model)
}

func (t *fileVisitor) VisitStructType(_ string, _ *ast.StructType) {

}

func (t *fileVisitor) walkTypeSpecNode(spec *ast.TypeSpec) {
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

func (t *fileVisitor) walkComments(commentGroup *ast.CommentGroup) {
	for _, comment := range commentGroup.List {
		t.VisitComment(comment)
	}
}

func (t *fileVisitor) newSymbolId() uint64 {
	next := t.idSequence + 1
	t.idSequence = next
	return next
}
