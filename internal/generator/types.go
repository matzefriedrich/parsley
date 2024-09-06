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
	Results    []Parameter
}

func (m Method) HasResults() bool {
	return len(m.Results) > 0
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

func (m Method) FormattedResultParameters() string {
	formattedResults := make([]string, len(m.Results))
	for i, result := range m.Results {
		formattedResults[i] = fmt.Sprintf("%s", result.Name)
	}
	return strings.Join(formattedResults, ", ")
}

func (m Method) FormattedResultTypes() string {
	formattedResults := make([]string, len(m.Results))
	for i, result := range m.Results {
		formattedResults[i] = fmt.Sprintf("%s", result.TypeName)
	}
	if len(formattedResults) == 0 {
		return ""
	}
	return "(" + strings.Join(formattedResults, ", ") + ")"
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

type ModelConfigurationFunc func(m *Model)

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
