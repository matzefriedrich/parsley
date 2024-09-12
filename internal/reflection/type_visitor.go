package reflection

import (
	"fmt"
	"go/ast"
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
		paramType := param.Type
		paramArrayType, isArrayType := paramType.(*ast.ArrayType)
		if isArrayType {
			paramType = paramArrayType.Elt
		}
		paramTypeIdentifier, _ := paramType.(*ast.Ident)
		for _, paramName := range param.Names {
			parameters = append(parameters, Parameter{
				Name:     paramName.Name,
				TypeName: paramTypeIdentifier.Name,
				IsArray:  isArrayType,
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
		fieldType := field.Type
		fieldArrayType, isArrayType := fieldType.(*ast.ArrayType)
		if isArrayType {
			fieldType = fieldArrayType.Elt
		}
		fieldTypeIdentifier, ok := fieldType.(*ast.Ident)
		if ok {
			parameters = append(parameters, Parameter{
				Name:     fmt.Sprintf("result%d", index),
				TypeName: fieldTypeIdentifier.Name,
				IsArray:  isArrayType,
			})
		}
	}
	return parameters
}
