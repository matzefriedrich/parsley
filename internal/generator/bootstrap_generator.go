package generator

import (
	"github.com/matzefriedrich/parsley/internal/templates"
	"io"
	"path"
)

type bootstrapGenerator struct {
	writerFunc ScaffoldingFileWriterFunc
}

// BootstrapGenerator provides functionalities to scaffold project files.
type BootstrapGenerator interface {
	ScaffoldProjectFiles()
}

var _ BootstrapGenerator = (*bootstrapGenerator)(nil)

type projectTemplateModel struct {
}

// ProjectItem represents a template file and its corresponding target filename in a project scaffold.
type ProjectItem struct {
	TemplateName   string
	TargetFilename string
}

type ScaffoldingFileWriterFunc func(targetFilename string) (io.WriteCloser, error)

// ScaffoldProjectFiles generates required project files using predefined templates and saves them to the project folder
func (b *bootstrapGenerator) ScaffoldProjectFiles() {

	gen := NewGenericCodeGenerator(func(name string) (string, error) {
		templateFilePath := path.Join("bootstrap", name)
		data, err := templates.BootstrapTemplates.ReadFile(templateFilePath)
		if err != nil {
			return "", err
		}
		return string(data), nil
	})

	m := &projectTemplateModel{}

	projectTemplateFiles := []ProjectItem{
		{TemplateName: "application.gotmpl", TargetFilename: "application.go"},
		{TemplateName: "main.gotmpl", TargetFilename: "main.go"},
		{TemplateName: "greeter.gotmpl", TargetFilename: "greeter.go"},
	}

	for _, projectItem := range projectTemplateFiles {
		generateFile := func(item ProjectItem) error {
			f, err := b.writerFunc(item.TargetFilename)
			if err != nil {
				return err
			}
			defer func(f io.WriteCloser) {
				_ = f.Close()
			}(f)
			return gen.Generate(item.TemplateName, m, f)
		}
		err := generateFile(projectItem)
		if err != nil {
			continue
		}
	}
}

// NewBootstrapGenerator initializes and returns a BootstrapGenerator with the specified project folder path.
func NewBootstrapGenerator(writerFunc ScaffoldingFileWriterFunc) BootstrapGenerator {
	return &bootstrapGenerator{
		writerFunc: writerFunc,
	}
}
