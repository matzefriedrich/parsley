package generator

import (
	"go/ast"
	"go/parser"
	"go/token"
)

type AstFileAccessor func() (*ast.File, error)

// AstFromFile Creates an AstFileAccessor object for the given Golang source file.
func AstFromFile(sourceFilePath string) AstFileAccessor {
	return func() (*ast.File, error) {
		fileSet := token.NewFileSet()
		return parser.ParseFile(fileSet, sourceFilePath, nil, parser.ParseComments)
	}
}

// AstFromSource Creates an AstFileAccessor object for the given source code.
func AstFromSource(code []byte) AstFileAccessor {
	return func() (*ast.File, error) {
		fileSet := token.NewFileSet()
		return parser.ParseFile(fileSet, "", code, parser.ParseComments)
	}
}
