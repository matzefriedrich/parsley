package generator

import (
	"fmt"
	"go/ast"
)

type TemplateModelBuilder struct {
	node *ast.File
}

func NewTemplateModelBuilder(accessor AstFileAccessor) (*TemplateModelBuilder, error) {
	node, err := accessor()
	if err != nil {
		return nil, err
	}
	return &TemplateModelBuilder{
		node: node,
	}, nil
}

func (b *TemplateModelBuilder) Build() (*Model, error) {
	packageName := b.node.Name.Name
	m := NewModel(packageName)
	b.collectImports(m)
	b.collectInterfaces(m)
	return m, nil
}

func (b *TemplateModelBuilder) collectImports(m *Model) {
	ast.Inspect(b.node, func(n ast.Node) bool {
		importSpec, ok := n.(*ast.ImportSpec)
		if ok {
			m.Imports = append(m.Imports, importSpec.Path.Value)
		}
		return true
	})
}

func (b *TemplateModelBuilder) collectInterfaces(m *Model) {
	ast.Inspect(b.node, func(n ast.Node) bool {
		typeSpec, interfaceType, ok := isInterfaceType(n)
		if ok {
			interfaceModel := InterfaceWithName(typeSpec.Name.Name)
			b.collectMethodsFor(interfaceType, &interfaceModel)
			m.Interfaces = append(m.Interfaces, interfaceModel)
		}
		return true
	})
}

func (b *TemplateModelBuilder) collectMethodsFor(interfaceType *ast.InterfaceType, interfaceModel *Interface) {
	for _, method := range interfaceType.Methods.List {
		if funcType, ok := method.Type.(*ast.FuncType); ok {
			name := method.Names[0].Name
			parameters := b.collectParametersFor(funcType)
			results := b.collectResultFieldsFor(funcType)
			interfaceModel.Methods = append(interfaceModel.Methods, Method{
				Name:       name,
				Parameters: parameters,
				Results:    results,
			})
		}
	}
}

func (b *TemplateModelBuilder) collectParametersFor(funcType *ast.FuncType) []Parameter {
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

func (b *TemplateModelBuilder) collectResultFieldsFor(funcType *ast.FuncType) []Parameter {
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

func isInterfaceType(n ast.Node) (*ast.TypeSpec, *ast.InterfaceType, bool) {
	typeSpec, ok := n.(*ast.TypeSpec)
	if ok {
		interfaceType, isInterface := typeSpec.Type.(*ast.InterfaceType)
		return typeSpec, interfaceType, isInterface
	}
	return nil, nil, false
}
