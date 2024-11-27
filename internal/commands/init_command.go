package commands

import (
	"fmt"
	"github.com/matzefriedrich/cobra-extensions/pkg/commands"
	"github.com/matzefriedrich/cobra-extensions/pkg/types"
	"github.com/matzefriedrich/parsley/internal/utils"
	"io"
	"os"
	"path"

	"github.com/matzefriedrich/parsley/internal/generator"
	"github.com/spf13/cobra"
)

// ScaffoldingFileWriterFactoryFunc defines a function type that returns a generator.ScaffoldingFileWriterFunc.
type ScaffoldingFileWriterFactoryFunc func(projectFolder string) (generator.ScaffoldingFileWriterFunc, error)

// ProjectLoaderFunc is a function type that loads a Go project given a project folder path.
type ProjectLoaderFunc func(projectFolderPath string) (generator.GoProject, error)

type initCommand struct {
	use                   types.CommandName `flag:"init" short:"Add Parsley to an application" long:"Integrates Parsley into an existing application by setting up the necessary scaffolding for dependency injection and code generation. It initializes project configurations, generates essential files, and prepares the application for using Parsley's advanced features."`
	fileWriterFactoryFunc ScaffoldingFileWriterFactoryFunc
	projectLoadFunc       ProjectLoaderFunc
}

// Execute sets up a new project by loading the current project folder, adding necessary dependencies,
// and generating initial project files. It handles errors and ensures the minimum required version of dependencies.
func (g *initCommand) Execute() {

	projectFolderPath, _ := os.Getwd()
	p, err := g.projectLoadFunc(projectFolderPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	minVersion, _ := utils.ApplicationVersion()
	const packageName = "github.com/matzefriedrich/parsley"
	dependencyErr := p.AddDependency(packageName, minVersion.String())
	if dependencyErr != nil {
		fmt.Println(err)
		return
	}

	fileWriterFunc, _ := g.fileWriterFactoryFunc(projectFolderPath)

	gen := generator.NewBootstrapGenerator(fileWriterFunc)
	gen.ScaffoldProjectFiles()
}

var _ types.TypedCommand = &initCommand{}

func NewInitCommand(
	writerFactoryFunc ScaffoldingFileWriterFactoryFunc,
	projectLoaderFunc ProjectLoaderFunc) *cobra.Command {
	command := &initCommand{
		fileWriterFactoryFunc: writerFactoryFunc,
		projectLoadFunc:       projectLoaderFunc,
	}
	return commands.CreateTypedCommand(command)
}

// NewProjectFileScaffoldingWriterFactory creates a factory for generating file writers in a specified project directory.
func NewProjectFileScaffoldingWriterFactory(projectFolderPath string) generator.ScaffoldingFileWriterFunc {
	return func(targetFilename string) (io.WriteCloser, error) {
		targetFilePath := path.Join(projectFolderPath, targetFilename)
		f, fileErr := os.OpenFile(targetFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if fileErr != nil {
			return nil, fileErr
		}
		return f, nil
	}
}

// LoadProjectFromDisk loads a Go project from the specified directory path.
func LoadProjectFromDisk(projectFolderPath string) (generator.GoProject, error) {
	return generator.OpenProject(projectFolderPath)
}
