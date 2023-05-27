// Package make implements the command to generate a new file
package make

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/awgst/gig/pkg"
	content "github.com/awgst/gig/template/request"
	"github.com/spf13/cobra"
)

// requestData is the data that will be parsed in the template
type requestData struct {
	Name string
}

// requestOptions is the options that will be parsed in the command
type RequestOptions struct {
	Module string
}

var request RequestOptions

// RequestCommand is the command to generate a new request
var RequestCommand = &cobra.Command{
	Use:   "make:request <name>",
	Short: "Make a new request",
	Long:  "Make a new request",
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
	Run: runMakeRequestCommand,
}

func init() {
	flags := RequestCommand.Flags()
	flags.StringVar(&request.Module, "module", "", "Specify the module that will be the destination")
}

// runMakeRequestCommand is the function that will be executed when the command is called
// It will generate a new request file based on the template
// The template will be parsed with the requestData
// Accepts the name of the module as an argument
func runMakeRequestCommand(cmd *cobra.Command, args []string) {
	GenerateRequest(request, args)
}

// GenerateRequest is the function that will be executed when the command is called
// It will generate a new request file based on the template
// The template will be parsed with the requestData
// Accepts the name of the module as an argument
func GenerateRequest(requestOpt RequestOptions, args []string) {
	templateContent := content.RequestTemplate
	name := strings.ToLower(args[0])
	moduleName := name
	if requestOpt.Module != "" {
		moduleName = requestOpt.Module
	}
	fileName := fmt.Sprintf("%s_request", name)
	requestData := requestData{
		Name: pkg.SnakeToPascal(name),
	}

	// Generate request file based on template
	pkg.GenerateFile("http/request", fileName, moduleName, templateContent, requestData)
}
