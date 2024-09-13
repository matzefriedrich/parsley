package commands

import (
	"fmt"
	"github.com/matzefriedrich/cobra-extensions/pkg"
	"github.com/matzefriedrich/cobra-extensions/pkg/abstractions"
	"github.com/matzefriedrich/parsley/internal/generator"
	"github.com/matzefriedrich/parsley/internal/reflection"
	"github.com/matzefriedrich/parsley/internal/templates"
	"github.com/spf13/cobra"
	"slices"
	"strings"
)

type mocksGeneratorCommand struct {
	use abstractions.CommandName `flag:"mocks" short:"Generate configurable mocks for interface types."`
}

type MocksGeneratorBehavior int

const (
	Default MocksGeneratorBehavior = iota
	OnlyMarked
	ExcludeIgnored
)

type ParsleyMockAnnotationAttribute int

const (
	Mock ParsleyMockAnnotationAttribute = iota + 1
	Ignore
)

func (p ParsleyMockAnnotationAttribute) String() string {
	switch p {
	case Mock:
		return "mock"
	case Ignore:
		return "ignore"
	default:
		return ""
	}
}

func (m *mocksGeneratorCommand) Execute() {

	templateLoader := func(_ string) (string, error) {
		return templates.MockTemplate, nil
	}

	kind := "mocks"
	gen, _ := generator.NewCodeFileGenerator(kind, func(config *generator.CodeFileGeneratorOptions) {
		config.TemplateLoader = templateLoader
		config.ConfigureModelCallback = func(m *reflection.Model) {
			m.AddImport("github.com/matzefriedrich/parsley/pkg/features")
			m.Interfaces = filterInterfaces(m)
		}
	})

	err := gen.GenerateCode()
	if err != nil {
		fmt.Println(err)
	}
}

func filterInterfaces(m *reflection.Model) []reflection.Interface {

	behavior := determineMockGeneratorBehavior(m)

	filterIdentifiers := func(attribute ParsleyMockAnnotationAttribute) []uint64 {
		identifiers := make([]uint64, 0)
		for _, comment := range m.Comments {
			if isParsleyMockDirective(comment, attribute) {
				p := comment.Pos
				for _, t := range m.Interfaces {
					if t.Pos > p {
						identifiers = append(identifiers, t.Id)
						break
					}
				}
			}
		}
		return identifiers
	}

	switch behavior {
	case OnlyMarked:
		keep := filterIdentifiers(Mock)
		// Keep interfaces whose identifier is in the keep slice
		return slices.DeleteFunc(m.Interfaces, func(i reflection.Interface) bool {
			return !slices.Contains(keep, i.Id)
		})
	case ExcludeIgnored:
		removed := filterIdentifiers(Ignore)
		// Remove interfaces whose identifier is in the removed slice
		return slices.DeleteFunc(m.Interfaces, func(i reflection.Interface) bool {
			return slices.Contains(removed, i.Id)
		})
	default:
		return m.Interfaces
	}
}

var _ pkg.TypedCommand = (*mocksGeneratorCommand)(nil)

func NewGenerateMocksCommand() *cobra.Command {
	command := &mocksGeneratorCommand{}
	return pkg.CreateTypedCommand(command)
}

func determineMockGeneratorBehavior(m *reflection.Model) MocksGeneratorBehavior {

	hasMockAnnotations := slices.ContainsFunc(m.Comments, func(comment reflection.Comment) bool {
		return isParsleyMockDirective(comment, Mock)
	})

	if hasMockAnnotations {
		return OnlyMarked
	}

	hasIgnoreAnnotations := slices.ContainsFunc(m.Comments, func(comment reflection.Comment) bool {
		return isParsleyMockDirective(comment, Ignore)
	})

	if hasIgnoreAnnotations {
		return ExcludeIgnored
	}

	return Default
}

// isGenerateMockDirective Returns true if the comment matches the directive, otherwise false.
func isParsleyMockDirective(comment reflection.Comment, annotation ParsleyMockAnnotationAttribute) bool {

	commentText := strings.TrimSpace(comment.Text)
	words := strings.Fields(commentText)

	annotationString := annotation.String()
	if annotationString == "" {
		return false
	}

	expected := []string{fmt.Sprintf("//parsley:%s", annotationString)}

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
