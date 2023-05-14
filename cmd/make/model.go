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

type modelData struct {
	Name string
}

type modelOptions struct {
	Module string
	Plain  bool
}

var model modelOptions

var ModelCommand = &cobra.Command{
	Use:   "make:model",
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
	model = modelOptions{}
	flags := ModelCommand.Flags()
	flags.StringVar(&model.Module, "module", "", "Specify the module that will be the destination")
	flags.BoolVar(&model.Plain, "plain", false, "Create a plain model struct without any field")
}

func runMakeModelCommand(cmd *cobra.Command, args []string) {
	fileName := strings.ToLower(args[0])
	templateContent := content.ModelTemplate
	moduleName := fileName
	if model.Module != "" {
		moduleName = model.Module
	}
	// Generate model file based on template
	modelData := modelData{
		Name: pkg.SnakeToPascal(args[0]),
	}

	if model.Plain {
		templateContent = content.PlainModelTemplate
	}
	pkg.GenerateFile("model", fileName, moduleName, templateContent, modelData)
}
