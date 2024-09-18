package main

import (
	"github.com/matzefriedrich/cobra-extensions/pkg/charmer"
	"github.com/matzefriedrich/parsley/internal/commands"
	"github.com/matzefriedrich/parsley/internal/generator"
)

func main() {

	description := "Welcome to Parsley \U0001F33F CLI! \n\n" +
		"Simplifying Advanced Dependency Injection in Go \n\n" +
		"With Parsley, you can generate boilerplate code for advanced DI features effortlessly. " +
		"Whether you're working with proxies, decorators, or need support for dynamic dependency resolution, " +
		"Parsley CLI has you covered. Focus on your core business logic while it takes care of the heavy lifting."

	app := charmer.NewCommandLineApplication("parsley-cli", description)

	app.AddCommand(
		commands.NewInitCommand(),
		commands.NewVersionCommand())

	app.AddGroupCommand(
		commands.NewGenerateGroupCommand(),
		func(w charmer.CommandSetup) {
			goFileAccessor := generator.GoFileAccessor()
			outputWriterFactory := generator.FileOutputWriter()
			w.AddCommand(commands.NewGenerateMocksCommand(goFileAccessor, outputWriterFactory))
			w.AddCommand(commands.NewGenerateProxyCommand())
		})

	_ = app.Execute()
}
