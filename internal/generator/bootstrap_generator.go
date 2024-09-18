package generator

import (
	"github.com/matzefriedrich/parsley/internal/templates"
	"os"
	"path"
)

type bootstrapGenerator struct {
	projectFolderPath string
}

type BootstrapGenerator interface {
	GenerateProjectFiles()
}

var _ BootstrapGenerator = (*bootstrapGenerator)(nil)

type ProjectTemplateModel struct {
}

type ProjectItem struct {
	TemplateName   string
	TargetFilename string
}

func (b *bootstrapGenerator) GenerateProjectFiles() {

	gen := NewGenericCodeGenerator(func(name string) (string, error) {
		templateFilePath := path.Join("bootstrap", name)
		data, err := templates.BootstrapTemplates.ReadFile(templateFilePath)
		if err != nil {
			return "", err
		}
		return string(data), nil
	})

	m := &ProjectTemplateModel{}

	projectTemplateFiles := []ProjectItem{
		{TemplateName: "application.gotmpl", TargetFilename: "application.go"},
		{TemplateName: "main.gotmpl", TargetFilename: "main.go"},
		{TemplateName: "greeter.gotmpl", TargetFilename: "greeter.go"},
	}

	for _, v := range projectTemplateFiles {
		generateFile := func(item ProjectItem) error {
			targetFilePath := path.Join(b.projectFolderPath, item.TargetFilename)
			f, _ := os.OpenFile(targetFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
			defer func(f *os.File) {
				_ = f.Close()
			}(f)
			return gen.Generate(item.TemplateName, m, f)
		}
		err := generateFile(v)
		if err != nil {
			continue
		}
	}
}

func NewBootstrapGenerator(projectFolderPath string) BootstrapGenerator {
	return &bootstrapGenerator{
		projectFolderPath: projectFolderPath,
	}
}
