package reflection

import (
	"github.com/matzefriedrich/parsley/internal/reflection"
	"github.com/matzefriedrich/parsley/pkg/features"
	"go/ast"
	"testing"
)

func Test_typeWalker_WalkInterface_skips_unnamed_methods(t *testing.T) {

	// Arrange
	typeVisitor := newAstTypeVisitorMock()

	sut := reflection.NewTypeWalker(typeVisitor)
	interfaceType := &ast.InterfaceType{
		Methods: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{},
				},
			}}}

	// Act
	sut.WalkInterface(interfaceType)

	// Assert
	typeVisitor.Verify(functionVisitMethodName, features.TimesNever())
}

func Test_typeWalker_WalkInterface_invokes_VisitMethod(t *testing.T) {

	// Arrange
	typeVisitor := newAstTypeVisitorMock()

	sut := reflection.NewTypeWalker(typeVisitor)

	expectedMethodName := "Method0"
	interfaceType := &ast.InterfaceType{
		Methods: &ast.FieldList{
			List: []*ast.Field{
				{
					Type: &ast.FuncType{},
					Names: []*ast.Ident{
						{
							Name: expectedMethodName,
						},
					},
				},
			}}}

	// Act
	sut.WalkInterface(interfaceType)

	// Assert
	typeVisitor.Verify(functionVisitMethodName, features.TimesOnce(), features.Exact(expectedMethodName), features.IsAny())
}

const (
	functionVisitMethodName = "VisitMethod"
)

type visitMethodFunc func(name string, funcType *ast.FuncType)

type astTypeVisitorMock struct {
	features.MockBase
	visitMethodFunc visitMethodFunc
}

func (a *astTypeVisitorMock) VisitMethod(name string, funcType *ast.FuncType) {
	a.MockBase.TraceMethodCall(functionVisitMethodName, name, funcType)
	a.visitMethodFunc(name, funcType)
}

var _ reflection.AstTypeVisitor = (*astTypeVisitorMock)(nil)

func newAstTypeVisitorMock() *astTypeVisitorMock {
	mock := &astTypeVisitorMock{
		MockBase: features.NewMockBase(),
		visitMethodFunc: func(name string, funcType *ast.FuncType) {
		},
	}
	mock.AddFunction(functionVisitMethodName, "VisitMethod(name string, funcType *ast.FuncType)")
	return mock
}
