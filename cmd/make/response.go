// Package make implements the command to generate a new file
package make

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/awgst/gig/pkg"
	content "github.com/awgst/gig/template/response"
	"github.com/spf13/cobra"
)

// responseData is the data that will be parsed in the template
type responseData struct {
	Name          string
	CamelCaseName string
	ModelName     string
	ModelPath     string
}

// responseOptions is the options that will be parsed in the command
type ResponseOptions struct {
	Module string
	Plain  bool
}

var response ResponseOptions

// ResponseCommand is the command to generate a new response
var ResponseCommand = &cobra.Command{
	Use:   "make:response <name>",
	Short: "Make a new response",
	Long:  "Make a new response",
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
	Run: runMakeResponseCommand,
}

func init() {
	flags := ResponseCommand.Flags()
	flags.StringVar(&response.Module, "module", "", "Specify the module that will be the destination")
	flags.BoolVar(&response.Plain, "plain", false, "Create a plain response struct without any field")
}

// runMakeResponseCommand is the function that will be executed when the command is called
// It will generate a new response file based on the template
// The template will be parsed with the responseData
// Accepts the name of the module as an argument
func runMakeResponseCommand(cmd *cobra.Command, args []string) {
	GenerateResponse(response, args)
}

// GenerateResponse is the function that will be executed when the command is called
// It will generate a new response file based on the template
// The template will be parsed with the responseData
// Accepts the name of the module as an argument
func GenerateResponse(responseOpt ResponseOptions, args []string) {
	templateContent := content.ResponseTemplate
	name := strings.ToLower(args[0])
	moduleName := name
	if response.Module != "" {
		moduleName = response.Module
	}
	fileName := fmt.Sprintf("%s_response", name)
	path := filepath.Join(projectName, "src/app", moduleName)
	responseData := responseData{
		Name:          pkg.SnakeToPascal(name),
		CamelCaseName: pkg.SnakeToCamel(name),
		ModelName:     pkg.SnakeToPascal(name),
		ModelPath:     path,
	}

	if responseOpt.Plain {
		templateContent = content.PlainResponseTemplate
	}

	// Generate repository file based on template
	pkg.GenerateFile("http/response", fileName, moduleName, templateContent, responseData)
}
