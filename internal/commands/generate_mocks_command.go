package commands

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/matzefriedrich/cobra-extensions/pkg/commands"
	"github.com/matzefriedrich/cobra-extensions/pkg/types"
	"github.com/matzefriedrich/parsley/internal/generator"
	"github.com/matzefriedrich/parsley/internal/reflection"
	"github.com/matzefriedrich/parsley/internal/templates"
	"github.com/spf13/cobra"
)

type mocksGeneratorCommand struct {
	types.BaseCommand   `cobra-x:"mocks, help='Generate configurable mocks for interface types.', description='Generates fully configurable mock implementations for Go interface types. It simplifies the process of creating mocks by analyzing the source code and automatically generating mock structs that adhere to the defined interfaces.'"`
	fileAccessor        reflection.AstFileAccessor
	outputWriterFactory generator.OutputWriterFactory
}

// MocksGeneratorBehavior defines different behaviors for mock generation, which influence how mocks are handled based on annotations.
type MocksGeneratorBehavior int

const (
	Default MocksGeneratorBehavior = iota
	OnlyMarked
	ExcludeIgnored
)

// ParsleyMockAnnotationAttribute represents an attribute used to manage the behavior of mocking annotations during code generation, such as including or ignoring specified interface mocks based on annotations.
type ParsleyMockAnnotationAttribute int

const (
	Mock ParsleyMockAnnotationAttribute = iota + 1
	Ignore
)

// String provides a string representation of the ParsleyMockAnnotationAttribute enum.
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

// Execute generates configurable mock implementations for all relevant interface types in the input source.
// It loads the template, configures the model, and writes the generated mocks to the specified output. The method also processes any errors encountered during code generation.
func (m *mocksGeneratorCommand) Execute(_ context.Context) {

	templateLoader := func(_ string) (string, error) {
		return templates.MockTemplate, nil
	}

	kind := "mocks"
	gen, _ := generator.NewCodeFileGenerator(kind, m.fileAccessor, func(config *generator.CodeFileGeneratorOptions) {
		config.TemplateLoader = templateLoader
		config.OutputWriterFactory = m.outputWriterFactory
		config.ConfigureModelCallback = func(m *reflection.Model) {
			m.Interfaces = filterInterfaces(m)
			m.Imports = mockImports(m)
		}
	})

	err := gen.GenerateCode()
	for err != nil {
		fmt.Printf("%+v\n", err)
		err = errors.Unwrap(err)
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

func mockImports(model *reflection.Model) []string {
	if len(model.Interfaces) == 0 {
		return nil
	}
	usedPaths := referencedImportPaths(model)
	usedPaths["github.com/matzefriedrich/parsley/pkg/features"] = struct{}{}
	return importsInSourceOrder(model.Imports, usedPaths)
}

func referencedImportPaths(model *reflection.Model) map[string]struct{} {
	pathsByQualifier := importPathsByQualifier(model)
	usedPaths := make(map[string]struct{})
	for _, interfaceModel := range model.Interfaces {
		for _, method := range interfaceModel.Methods {
			addParameterImports(method.Parameters, pathsByQualifier, usedPaths)
			addParameterImports(method.Results, pathsByQualifier, usedPaths)
		}
	}
	return usedPaths
}

func importPathsByQualifier(model *reflection.Model) map[string]string {
	pathsByQualifier := make(map[string]string, len(model.Imports))
	for _, importPath := range model.Imports {
		qualifier := model.ImportAliases[importPath]
		if qualifier == "" {
			parts := strings.Split(importPath, "/")
			qualifier = parts[len(parts)-1]
		}
		pathsByQualifier[qualifier] = importPath
	}
	return pathsByQualifier
}

func addParameterImports(parameters []reflection.Parameter, pathsByQualifier map[string]string, usedPaths map[string]struct{}) {
	for _, parameter := range parameters {
		for typePart := parameter.Type; typePart != nil; typePart = typePart.Next {
			if importPath, found := pathsByQualifier[typePart.SelectorName]; found {
				usedPaths[importPath] = struct{}{}
			}
		}
	}
}

func importsInSourceOrder(sourceImports []string, usedPaths map[string]struct{}) []string {
	imports := make([]string, 0, len(usedPaths))
	for _, importPath := range sourceImports {
		if _, used := usedPaths[importPath]; used {
			imports = append(imports, importPath)
		}
	}
	featuresPath := "github.com/matzefriedrich/parsley/pkg/features"
	if !slices.Contains(imports, featuresPath) {
		imports = append(imports, featuresPath)
	}
	return imports
}

var _ types.TypedCommand = (*mocksGeneratorCommand)(nil)

// NewGenerateMocksCommand creates a new cobra command to generate mock implementations for interfaces.
// This command uses the provided file accessor to read the source file and the output writer factory to write the generated mocks.
func NewGenerateMocksCommand(fileAccessor reflection.AstFileAccessor, outputWriterFactory generator.OutputWriterFactory) *cobra.Command {
	if fileAccessor == nil {
		panic("file accessor required")
	}
	if outputWriterFactory == nil {
		panic("output writer factory required")
	}
	command := &mocksGeneratorCommand{
		fileAccessor:        fileAccessor,
		outputWriterFactory: outputWriterFactory,
	}
	return commands.CreateTypedCommand(command)
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
