package cmd

import (
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

Help faster Go project development
More informations --> https://github.com/awgst/gig`,
}

func Execute() {
	rootCommand.CompletionOptions.DisableDefaultCmd = true
	_ = rootCommand.Execute()
}
