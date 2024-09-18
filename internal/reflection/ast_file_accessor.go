package reflection

import (
	"go/ast"
	"go/parser"
	"go/token"
)

type AstFileSource struct {
	File     *ast.File
	Filename string
}

type AstFileAccessor func() (*AstFileSource, error)

// AstFromFile Creates an AstFileAccessor object for the given Golang source file.
func AstFromFile(sourceFilePath string) AstFileAccessor {
	return func() (*AstFileSource, error) {
		fileSet := token.NewFileSet()
		f, err := parser.ParseFile(fileSet, sourceFilePath, nil, parser.ParseComments)
		source := &AstFileSource{File: f, Filename: sourceFilePath}
		return source, err
	}
}

// AstFromSource Creates an AstFileAccessor object for the given source code.
func AstFromSource(code []byte) AstFileAccessor {
	return func() (*AstFileSource, error) {
		fileSet := token.NewFileSet()
		const filename = ""
		f, err := parser.ParseFile(fileSet, filename, code, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		source := &AstFileSource{File: f, Filename: filename}
		return source, err
	}
}
