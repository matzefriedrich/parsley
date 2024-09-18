package reflection

import (
	"fmt"
	"github.com/matzefriedrich/parsley/internal"
	"go/ast"
)

type fieldTypeInfo struct {
	Name      string
	IsArray   bool
	IsPointer bool
}

type interfaceMethodsCollector struct {
	model *Interface
}

func (t *interfaceMethodsCollector) VisitMethod(name string, method *ast.FuncType) {
	parameters := CollectParametersFor(method)
	results := CollectResultFieldsFor(method)
	t.model.Methods = append(t.model.Methods, Method{
		Name:       name,
		Parameters: parameters,
		Results:    results,
	})
}

var _ AstTypeVisitor = (*interfaceMethodsCollector)(nil)

func newInterfaceMethodsCollector(interfaceModel *Interface) AstTypeVisitor {
	return &interfaceMethodsCollector{
		model: interfaceModel,
	}
}

func CollectParametersFor(funcType *ast.FuncType) []Parameter {
	parameters := make([]Parameter, 0)
	for _, param := range funcType.Params.List {
		typeInfo := getFieldTypeInfo(param)
		for _, paramName := range param.Names {
			parameters = append(parameters, Parameter{
				Name: paramName.Name,
				Type: typeInfo,
			})
		}
	}
	return parameters
}

func CollectResultFieldsFor(funcType *ast.FuncType) []Parameter {
	parameters := make([]Parameter, 0)
	if funcType.Results == nil {
		return parameters
	}
	for index, field := range funcType.Results.List {
		typeInfo := getFieldTypeInfo(field)
		parameters = append(parameters, Parameter{
			Name: fmt.Sprintf("result%d", index),
			Type: typeInfo,
		})
	}
	return parameters
}

func getFieldTypeInfo(param *ast.Field) *ParameterType {

	paramTypeName := ""

	paramType := param.Type

	typeStack := internal.MakeStack[ParameterType]()
	expressionStack := internal.MakeStack[ast.Expr]()
	expressionStack.Push(paramType)

	for expressionStack.IsEmpty() == false {

		next := expressionStack.Pop()

		switch next.(type) {
		case *ast.Ident:
			ident, _ := next.(*ast.Ident)
			paramTypeName = ident.Name
			typeStack.Push(ParameterType{Name: paramTypeName})

		case *ast.SelectorExpr:
			selector, _ := next.(*ast.SelectorExpr)
			ident, _ := selector.X.(*ast.Ident)
			typeStack.Push(ParameterType{SelectorName: ident.Name, Name: selector.Sel.Name})

		case *ast.ArrayType:
			arrayType := next.(*ast.ArrayType)
			t := arrayType.Elt
			typeStack.Push(ParameterType{IsArray: true})
			expressionStack.Push(t)

		case *ast.StarExpr:
			starExpr := next.(*ast.StarExpr)
			typeStack.Push(ParameterType{IsPointer: true})
			switch starExpr.X.(type) {
			case *ast.Ident:
				expressionStack.Push(starExpr.X)

			case *ast.SelectorExpr:
				selectorExpr, _ := starExpr.X.(*ast.SelectorExpr)
				expressionStack.Push(selectorExpr)

			case *ast.ArrayType:
				arrayType := starExpr.X.(*ast.ArrayType)
				expressionStack.Push(arrayType)
			}
		}
	}

	if typeStack.IsEmpty() {
		return nil
	}

	last := typeStack.Pop()
	result := &last
	for typeStack.IsEmpty() == false {
		parameterType := typeStack.Pop()
		parameterType.Next = result
		result = &parameterType
	}

	return result
}
