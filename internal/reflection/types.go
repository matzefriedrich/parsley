package reflection

import (
	"go/token"
	"strings"
)

type Parameter struct {
	Name string
	Type *ParameterType
}

// IsEllipsis determines if the parameter type is an ellipsis (`...`) parameter.
func (p Parameter) IsEllipsis() bool {
	pt := p.Type
	for pt != nil {
		if pt.IsEllipsis {
			return true
		}
		pt = pt.Next
	}
	return false
}

type ParameterType struct {
	Name         string
	SelectorName string
	IsArray      bool
	IsEllipsis   bool
	IsInterface  bool
	IsPointer    bool
	Next         *ParameterType
}

func (p Parameter) MatchesType(name string) bool {
	return strings.Compare(name, p.Type.Name) == 0
}

type SymbolInfo struct {
	Id  uint64
	Pos token.Pos
	End token.Pos
}

type Method struct {
	SymbolInfo
	Name       string
	Parameters []Parameter
	Results    []Parameter
}

type Interface struct {
	SymbolInfo
	Name    string
	Methods []Method
}

func InterfaceWithName(name string, info SymbolInfo) Interface {
	return Interface{
		SymbolInfo: info,
		Name:       name,
		Methods:    make([]Method, 0),
	}
}

type FuncType struct {
	SymbolInfo
	Name       string
	Parameters []Parameter
	Results    []Parameter
}

type Comment struct {
	SymbolInfo
	Text string
}

// Model The generator root model type.
type Model struct {
	Comments    []Comment
	Interfaces  []Interface
	FuncTypes   []FuncType
	PackageName string
	Imports     []string
}

type ModelConfigurationFunc func(m *Model)

func (m *Model) AddImport(s string) {
	s = strings.TrimSuffix(strings.TrimPrefix(s, "\""), "\"")
	m.Imports = append(m.Imports, s)
}
