package generator

import (
	"github.com/matzefriedrich/parsley/internal/reflection"
	"io"
)

type OutputWriterFactory func(kind string, source *reflection.AstFileSource) (io.WriteCloser, error)
