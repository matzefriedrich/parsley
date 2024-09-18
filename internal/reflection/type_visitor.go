package reflection

import (
	"fmt"
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
				Name:      paramName.Name,
				TypeName:  typeInfo.Name,
				IsArray:   typeInfo.IsArray,
				IsPointer: typeInfo.IsPointer,
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
			Name:      fmt.Sprintf("result%d", index),
			TypeName:  typeInfo.Name,
			IsArray:   typeInfo.IsArray,
			IsPointer: typeInfo.IsPointer,
		})
	}
	return parameters
}

func getFieldTypeInfo(param *ast.Field) fieldTypeInfo {

	paramTypeName := ""

	paramType := param.Type
	paramArrayType, isArrayType := paramType.(*ast.ArrayType)
	if isArrayType {
		paramType = paramArrayType.Elt
	}

	if paramType != nil {
		id, ok := paramType.(*ast.Ident)
		if ok {
			paramTypeName = id.Name
		}
	}

	paramPointerType, isPointerType := paramType.(*ast.StarExpr)
	if isPointerType && paramPointerType != nil {
		switch paramPointerType.X.(type) {
		case *ast.SelectorExpr:
			selectorExpr, _ := paramPointerType.X.(*ast.SelectorExpr)
			ident, xOk := selectorExpr.X.(*ast.Ident)
			if xOk {
				paramTypeName = fmt.Sprintf("%s.%s", ident.Name, selectorExpr.Sel.Name)
			}
			break
		}
	}

	return fieldTypeInfo{
		Name:      paramTypeName,
		IsArray:   isArrayType,
		IsPointer: isPointerType,
	}
}
