// Package make implements the command to generate a new file
package make

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/awgst/gig/pkg"
	content "github.com/awgst/gig/template/service"
	"github.com/spf13/cobra"
)

// serviceData is the data that will be parsed in the template
type serviceData struct {
	Name           string
	CamelCaseName  string
	LowerName      string
	ModelName      string
	ModelPath      string
	RequestPath    string
	RepositoryPath string
	RepositoryName string
}

// serviceOptions is the options that will be parsed in the command
type ServiceOptions struct {
	Module string
	CRUD   bool
}

var service ServiceOptions

// ServiceCommand is the command to generate a new service
var ServiceCommand = &cobra.Command{
	Use:   "make:service <name>",
	Short: "Make a new service",
	Long:  "Make a new service",
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
	Run: runMakeServiceCommand,
}

func init() {
	flags := ServiceCommand.Flags()
	flags.StringVar(&service.Module, "module", "", "Specify the module that will be the destination")
}

// runMakeServiceCommand is the function that will be executed when the command is called
// It will generate a new service file based on the template
// The template will be parsed with the serviceData
// Accepts the name of the module as an argument
func runMakeServiceCommand(cmd *cobra.Command, args []string) {
	GenerateService(service, args)
}

// GenerateService is the function that will be executed when the command is called
// It will generate a new service file based on the template
// The template will be parsed with the serviceData
// Accepts the name of the module as an argument
func GenerateService(serviceOpt ServiceOptions, args []string) {
	templateContent := content.ServiceTemplate
	name := strings.ToLower(args[0])
	moduleName := name
	if service.Module != "" {
		moduleName = service.Module
	}
	fileName := fmt.Sprintf("%s_service", name)
	path := filepath.Join(projectName, "src/app", moduleName)
	serviceData := serviceData{
		Name:           pkg.SnakeToPascal(name),
		LowerName:      strings.ToLower(name),
		CamelCaseName:  pkg.SnakeToCamel(name),
		ModelName:      pkg.SnakeToPascal(name),
		ModelPath:      path,
		RequestPath:    filepath.Join(path, "http"),
		RepositoryName: pkg.SnakeToPascal(name),
		RepositoryPath: path,
	}

	if serviceOpt.CRUD {
		templateContent = content.ServiceCRUDTemplate
	}

	// Generate repository file based on template
	pkg.GenerateFile("service", fileName, moduleName, templateContent, serviceData)
}
