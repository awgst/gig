// Package make implements the command to generate a new file
package make

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/awgst/gig/pkg"
	content "github.com/awgst/gig/template/handler"
	"github.com/spf13/cobra"
)

// handlerData is the data that will be parsed in the template
type handlerData struct {
	Name          string
	CamelCaseName string
	LowerName     string
	ModelName     string
	ModelPath     string
	RequestPath   string
	ResponsePath  string
	ServicePath   string
}

// handlerOptions is the options that will be parsed in the command
type HandlerOptions struct {
	Module string
	CRUD   bool
}

var handler HandlerOptions

// HandlerCommand is the command to generate a new handler
var HandlerCommand = &cobra.Command{
	Use:   "make:handler <name>",
	Short: "Make a new handler",
	Long:  "Make a new handler",
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
	Run: runMakeHandlerCommand,
}

func init() {
	flags := HandlerCommand.Flags()
	flags.StringVar(&handler.Module, "module", "", "Specify the module that will be the destination")
	flags.BoolVar(&handler.CRUD, "crud", false, "Create a handler with CRUD methods")
}

// runMakeHandlerCommand is the function that will be executed when the command is called
// It will generate a new handler file based on the template
// The template will be parsed with the handlerData
// Accepts the name of the module as an argument
func runMakeHandlerCommand(cmd *cobra.Command, args []string) {
	GenerateHandler(handler, args)
}

// GenerateHandler is the function that will be executed when the command is called
// It will generate a new handler file based on the template
// The template will be parsed with the handlerData
// Accepts the name of the module as an argument
func GenerateHandler(handlerOpt HandlerOptions, args []string) {
	templateContent := content.HandlerPlainTemplate
	name := strings.ToLower(args[0])
	moduleName := name
	if handlerOpt.Module != "" {
		moduleName = handlerOpt.Module
	}
	fileName := fmt.Sprintf("%s_handler", name)
	path := filepath.Join(projectName, "src/app", moduleName)
	httpPath := filepath.Join(path, "http")
	pascalName := pkg.SnakeToPascal(name)
	handlerData := handlerData{
		Name:          pascalName,
		LowerName:     strings.ToLower(name),
		CamelCaseName: pkg.SnakeToCamel(name),
		ModelName:     pascalName,
		ModelPath:     path,
		RequestPath:   httpPath,
		ResponsePath:  httpPath,
		ServicePath:   path,
	}

	httpFramework, _ := pkg.ReadJsonString("http_framework")
	if handlerOpt.CRUD {
		templateContent = content.HandlerCRUDTemplate[httpFramework]
	}

	// Generate handler file based on template
	pkg.GenerateFile("http/handler", fileName, moduleName, templateContent, handlerData)
}
