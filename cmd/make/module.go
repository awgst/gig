package make

import (
	"github.com/spf13/cobra"
)

var ModuleCommand = &cobra.Command{
	Use:   "make:module <name>",
	Short: "Make a new module",
	Long:  "Make a new module",
	Run:   runMakeModuleCommand,
}

func runMakeModuleCommand(cmd *cobra.Command, args []string) {
}
