// Package make implements the command to generate a new file
package make

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	content "github.com/awgst/gig/template/model"

	"github.com/awgst/gig/pkg"
	"github.com/spf13/cobra"
)

// modelData is the data that will be parsed in the template
type modelData struct {
	Name string
}

// ModelOptions is the options that will be parsed in the command
type ModelOptions struct {
	Module string
	Plain  bool
}

var model ModelOptions

// ModelCommand is the command to generate a new model
var ModelCommand = &cobra.Command{
	Use:   "make:model <name>",
	Short: "Make a new model",
	Long:  "Make a new model",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a name argument")
		}

		regex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
		if !regex.MatchString(args[0]) {
			return fmt.Errorf("invalid name: %s", args[0])
		}

		return nil
	},
	Run:     runMakeModelCommand,
	Example: `gig make:model user || gig make:model user_address`,
}

func init() {
	model = ModelOptions{}
	flags := ModelCommand.Flags()
	flags.StringVar(&model.Module, "module", "", "Specify the module that will be the destination")
	flags.BoolVar(&model.Plain, "plain", false, "Create a plain model struct without any field")
}

// runMakeModelCommand is the function that will be executed when the command is called
// It will generate a new model file based on the template
// The template will be parsed with the modelData
// Accepts the name of the module as an argument
func runMakeModelCommand(cmd *cobra.Command, args []string) {
	GenerateModel(model, args)
}

// GenerateModel is function that will be called from another package
// Accept ModelOptions and []string as an argument
func GenerateModel(modelOpt ModelOptions, args []string) {
	fileName := strings.ToLower(args[0])
	templateContent := content.ModelTemplate
	moduleName := fileName
	if modelOpt.Module != "" {
		moduleName = modelOpt.Module
	}
	// Generate model file based on template
	modelData := modelData{
		Name: pkg.SnakeToPascal(args[0]),
	}

	if modelOpt.Plain {
		templateContent = content.PlainModelTemplate
	}
	pkg.GenerateFile("model", fileName, moduleName, templateContent, modelData)
}
