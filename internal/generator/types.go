package generator

import (
	"fmt"
	"strings"
)

type Parameter struct {
	Name     string
	TypeName string
}

func (p Parameter) MatchesType(name string) bool {
	return strings.Compare(name, p.TypeName) == 0
}

type Method struct {
	Name       string
	Parameters []Parameter
}

func (m Method) FormattedParameters() string {
	formattedParameters := make([]string, len(m.Parameters))
	for i, parameter := range m.Parameters {
		formattedParameters[i] = fmt.Sprintf("%s %s", parameter.Name, parameter.TypeName)
	}
	return strings.Join(formattedParameters, ", ")
}

func (m Method) FormattedCallParameters() string {
	formattedParameters := make([]string, len(m.Parameters))
	for i, parameter := range m.Parameters {
		formattedParameters[i] = fmt.Sprintf("%s", parameter.Name)
	}
	return strings.Join(formattedParameters, ", ")
}

type Interface struct {
	Name    string
	Methods []Method
}

func InterfaceWithName(name string) Interface {
	return Interface{
		Name:    name,
		Methods: make([]Method, 0),
	}
}

// Model The generator root model type.
type Model struct {
	Interfaces  []Interface
	PackageName string
	Imports     []string
}

func (m *Model) AddImport(s string) {
	m.Imports = append(m.Imports, s)
}

func NewModel(packageName string) *Model {
	return &Model{
		PackageName: packageName,
		Interfaces:  make([]Interface, 0),
		Imports:     make([]string, 0),
	}
}
