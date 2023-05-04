package cmd

import (
	"errors"
	"fmt"
	"log"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

var createCommand = &cobra.Command{
	Use:     "create <name>",
	Aliases: []string{},
	Short:   "Create a new project",
	Long:    "Create a new project",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a name argument")
		}
		return nil
	},
	ValidArgs:  []string{"name"},
	ArgAliases: []string{"name"},
	Example:    "gig create go-project",
	Run:        runCreateCommand,
}

func init() {
	rootCommand.AddCommand(createCommand)
}

func runCreateCommand(cmd *cobra.Command, args []string) {
	err := survey.Ask(
		[]*survey.Question{
			{
				Name: "http_framework",
				Prompt: &survey.Select{
					Message: "Choose a http framework:",
					Options: []string{
						"chi",
						"echo",
						"fiber",
						"gin",
						"mux",
					},
					Default:  "chi",
					PageSize: 10,
				},
				Validate: survey.Required,
			},
		},
		&createCommandAnswer,
	)

	if err != nil {
		log.Fatal("Something went wrong : ", err)
	}

	fmt.Println(createCommandAnswer.HttpFramework, args[0])
}
