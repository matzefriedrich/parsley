package reflection

import (
	"fmt"
	"go/ast"

	"github.com/matzefriedrich/parsley/internal"
)

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
	resultIndex := 0
	for _, field := range funcType.Results.List {
		typeInfo := getFieldTypeInfo(field)
		if len(field.Names) == 0 {
			parameters = append(parameters, Parameter{
				Name: fmt.Sprintf("result%d", resultIndex),
				Type: typeInfo,
			})
			resultIndex++
		} else {
			for _, name := range field.Names {
				parameters = append(parameters, Parameter{
					Name: name.Name,
					Type: typeInfo,
				})
				resultIndex++
			}
		}
	}
	return parameters
}

func getFieldTypeInfo(param *ast.Field) *ParameterType {

	paramTypeName := ""

	paramType := param.Type

	typeStack := internal.MakeStack[ParameterType]()
	expressionStack := internal.MakeStack[ast.Expr]()
	expressionStack.Push(paramType)

	for !expressionStack.IsEmpty() {

		next := expressionStack.Pop()

		switch expr := next.(type) {
		case *ast.Ellipsis:
			typeStack.Push(ParameterType{IsEllipsis: true})
			expressionStack.Push(expr.Elt)

		case *ast.Ident:
			paramTypeName = expr.Name
			typeStack.Push(ParameterType{Name: paramTypeName})

		case *ast.InterfaceType:
			typeStack.Push(ParameterType{IsInterface: true})

		case *ast.SelectorExpr:
			ident, _ := expr.X.(*ast.Ident)
			typeStack.Push(ParameterType{SelectorName: ident.Name, Name: expr.Sel.Name})

		case *ast.ArrayType:
			t := expr.Elt
			typeStack.Push(ParameterType{IsArray: true})
			expressionStack.Push(t)

		case *ast.StarExpr:
			typeStack.Push(ParameterType{IsPointer: true})
			switch starExpr := expr.X.(type) {
			case *ast.Ident:
				expressionStack.Push(starExpr)
			case *ast.SelectorExpr:
				expressionStack.Push(starExpr)
			case *ast.ArrayType:
				expressionStack.Push(starExpr)
			}
		}
	}

	if typeStack.IsEmpty() {
		return nil
	}

	last := typeStack.Pop()
	result := new(last)
	for !typeStack.IsEmpty() {
		parameterType := typeStack.Pop()
		parameterType.Next = result
		result = &parameterType
	}

	return result
}
