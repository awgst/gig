// Package make implements the command to generate a new file
package make

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/awgst/gig/pkg"
	content "github.com/awgst/gig/template/module"
	"github.com/spf13/cobra"
)

// ModuleOptions is the options that will be parsed in the command
type ModuleOptions struct {
	CRUD bool
}

// ModuleCommand is the command to generate a new module
var ModuleCommand = &cobra.Command{
	Use:   "make:module <name>",
	Short: "Make a new module",
	Long:  "Make a new module",
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
	Run: runMakeModuleCommand,
}

// moduleData is the data that will be parsed in the template
type moduleData struct {
	Name           string
	CamelCaseName  string
	ModulePath     string
	HandlerName    string
	ServiceName    string
	RepositoryName string
	PackageName    string
	ProjectName    string
}

var module ModuleOptions

func init() {
	flags := ModuleCommand.Flags()
	flags.BoolVar(&module.CRUD, "crud", false, "Create a module with CRUD methods")
}

// runMakeModuleCommand is the function that will be executed when the command is called
func runMakeModuleCommand(cmd *cobra.Command, args []string) {
	GenerateModule(module, args)
}

// GenerateModule is the function that will be executed when the command is called
// It will generate handler, model, repository, request, response, service files based on the template
func GenerateModule(moduleOpt ModuleOptions, args []string) {
	moduleName := args[0]
	crud := moduleOpt.CRUD
	if crud {
		// Generate request
		requestOpt := RequestOptions{
			Module: moduleName,
		}
		GenerateRequest(requestOpt, []string{fmt.Sprintf("create_%s", moduleName)})
		GenerateRequest(requestOpt, []string{fmt.Sprintf("update_%s", moduleName)})

		// Generate response
		responseOpt := ResponseOptions{
			Module: moduleName,
		}
		GenerateResponse(responseOpt, []string{moduleName})
	}

	// Generate model
	modelOpt := ModelOptions{
		Module: moduleName,
	}
	GenerateModel(modelOpt, []string{moduleName})

	// Generate repository
	repositoryOpt := RepositoryOptions{
		Module: moduleName,
		CRUD:   crud,
	}
	GenerateRepository(repositoryOpt, []string{moduleName})

	// Generate service
	serviceOpt := ServiceOptions{
		Module: moduleName,
		CRUD:   crud,
	}
	GenerateService(serviceOpt, []string{moduleName})

	// Generate handler
	handlerOpt := HandlerOptions{
		Module: moduleName,
		CRUD:   crud,
	}
	GenerateHandler(handlerOpt, []string{moduleName})

	projectName, _ := pkg.ReadJsonString("name")
	pascalName := pkg.SnakeToPascal(moduleName)
	// Generate module file
	moduleData := moduleData{
		Name:           pascalName,
		CamelCaseName:  pkg.SnakeToCamel(moduleName),
		ModulePath:     filepath.Join(projectName, "src/app", moduleName),
		HandlerName:    fmt.Sprintf("%sHandler", pascalName),
		ServiceName:    fmt.Sprintf("%sService", pascalName),
		RepositoryName: fmt.Sprintf("%sRepository", pascalName),
		PackageName:    strings.ReplaceAll(moduleName, "_", ""),
		ProjectName:    projectName,
	}

	// Generate module file
	fileName := moduleName
	modulTemplate := content.ModuleTemplate
	if crud {
		modulTemplate = content.ModuleCRUDTemplate
	}
	pkg.GenerateFile("", fileName, moduleName, modulTemplate, moduleData)
}
