package main

import (
	"github.com/matzefriedrich/cobra-extensions/pkg/charmer"
	"github.com/matzefriedrich/parsley/internal/commands"
)

func main() {

	app := charmer.NewCommandLineApplication()

	app.AddGroupCommand(
		commands.NewGenerateGroupCommand(),
		func(w charmer.CommandSetup) {
			w.AddCommand(commands.NewGenerateProxyCommand())
		})

	app.Execute()
}
