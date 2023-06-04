// Package cmd implements the list of commands that can be executed
package migrate

import (
	"github.com/spf13/cobra"
)

var RootCommand = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate files",
	Long:  `Migrate files with suffix .up.sql and .down.sql`,
}

func init() {
	RootCommand.CompletionOptions.DisableDefaultCmd = true
	RootCommand.AddCommand(
		upCommand,
	)
}
