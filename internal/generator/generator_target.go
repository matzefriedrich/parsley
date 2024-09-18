package generator

import (
	"github.com/matzefriedrich/parsley/internal/reflection"
	"io"
)

type OutputWriter interface {
	io.Writer
	io.Closer
}

type OutputWriterFactory func(kind string, source *reflection.AstFileSource) (OutputWriter, error)
