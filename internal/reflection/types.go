package reflection

import (
	"go/token"
	"strings"
)

type Parameter struct {
	Name     string
	TypeName string
	IsArray  bool
}

func (p Parameter) MatchesType(name string) bool {
	return strings.Compare(name, p.TypeName) == 0
}

type SymbolInfo struct {
	Pos token.Pos
	End token.Pos
}

func NewSymbolInfo(pos token.Pos, end token.Pos) SymbolInfo {
	return SymbolInfo{pos, end}
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
	m.Imports = append(m.Imports, s)
}
