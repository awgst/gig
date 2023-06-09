// Package cmd implements the list of commands that can be executed
package cmd

import (
	"github.com/awgst/gig/cmd/make"
	"github.com/awgst/gig/cmd/migrate"
	"github.com/awgst/gig/pkg"
	"github.com/spf13/cobra"
)

var (
	createCommandAnswer pkg.CreateCommandAnswer
)

var rootCommand = &cobra.Command{
	Use:     "gig",
	Version: pkg.CLIVersion,
	Short:   "CLI for Go project",
	Long: `   ________________
  / ____/  _/ ____/
 / / __ / // / __  
/ /_/ // // /_/ /  
\____/___/\____/   

🚀 Command line interface which provides a number helpful command to assist Go project development!
❔ More informations --> https://github.com/awgst/gig`,
}

func Execute() {
	rootCommand.CompletionOptions.DisableDefaultCmd = true
	rootCommand.AddCommand(
		createCommand,
		upCommand,
		make.MigrationCommand,
		make.ModuleCommand,
		make.ModelCommand,
		make.RepositoryCommand,
		make.ServiceCommand,
		make.HandlerCommand,
		make.RequestCommand,
		make.ResponseCommand,
		migrate.RootCommand,
	)
	_ = rootCommand.Execute()
}
