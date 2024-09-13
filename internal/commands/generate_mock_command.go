package commands

import (
	"fmt"
	"github.com/matzefriedrich/cobra-extensions/pkg"
	"github.com/matzefriedrich/cobra-extensions/pkg/abstractions"
	"github.com/matzefriedrich/parsley/internal/generator"
	"github.com/matzefriedrich/parsley/internal/reflection"
	"github.com/matzefriedrich/parsley/internal/templates"
	"github.com/spf13/cobra"
	"strings"
)

type mockGeneratorCommand struct {
	use abstractions.CommandName `flag:"mock" short:"Generates configurable mocks for annotated interface types."`
}

func (m *mockGeneratorCommand) Execute() {

	templateLoader := func(_ string) (string, error) {
		return templates.MockTemplate, nil
	}

	kind := "mock"
	gen, _ := generator.NewCodeFileGenerator(kind, func(config *generator.CodeFileGeneratorOptions) {
		config.TemplateLoader = templateLoader
		config.ConfigureModelCallback = func(m *reflection.Model) {

			m.AddImport("github.com/matzefriedrich/parsley/pkg/features")

			interfaces := make([]reflection.Interface, 0)
			for _, comment := range m.Comments {
				if isGenerateMockDirective(comment) {
					p := comment.Pos
					for _, t := range m.Interfaces {
						if t.Pos > p {
							interfaces = append(interfaces, t)
							break
						}
					}
				}
			}
			m.Interfaces = interfaces
		}
	})

	err := gen.GenerateCode()
	if err != nil {
		fmt.Println(err)
	}
}

var _ pkg.TypedCommand = (*mockGeneratorCommand)(nil)

func NewGenerateMockCommand() *cobra.Command {
	command := &mockGeneratorCommand{}
	return pkg.CreateTypedCommand(command)
}

// isGenerateMockDirective Returns true if the comment matches the directive, otherwise false.
func isGenerateMockDirective(comment reflection.Comment) bool {

	commentText := strings.TrimSpace(comment.Text)
	words := strings.Fields(commentText)

	expected := []string{"//go:generate", "parsley-cli", "generate", "mock"}

	if len(words) >= len(expected) {
		for i := range expected {
			if words[i] != expected[i] {
				return false
			}
		}
		return true
	}
	return false
}
