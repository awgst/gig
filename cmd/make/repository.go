// Package make implements the command to generate a new file
package make

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/awgst/gig/pkg"
	content "github.com/awgst/gig/template/repository"
	"github.com/spf13/cobra"
)

// repositoryData is the data that will be parsed in the template
type repositoryData struct {
	Name      string
	LowerName string
	ModelPath string
}

// repositoryOptions is the options that will be parsed in the command
type repositoryOptions struct {
	CRUD   bool
	Module string
}

// repository is the options that will be parsed in the command
var repository repositoryOptions

// projectName is the name of the project
var projectName string

// RepositoryCommand is the command to generate a new repository
var RepositoryCommand = &cobra.Command{
	Use:   "make:repository",
	Short: "Make a new repository",
	Long:  "Make a new repository",
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
	Run: runMakeRepositoryCommand,
}

func init() {
	flags := RepositoryCommand.Flags()
	flags.BoolVar(&repository.CRUD, "crud", false, "Create a repository with CRUD methods")
	flags.StringVar(&repository.Module, "module", "", "Specify the module that will be the destination")
	projectName, _ = pkg.ReadJsonString("name")
}

// runMakeRepositoryCommand is the function that will be executed when the command is called
// It will generate a new repository file based on the template
// The template will be parsed with the repositoryData
// Accepts the name of the module as an argument
func runMakeRepositoryCommand(cmd *cobra.Command, args []string) {
	generateRepository(args)
	generateRepositoryInterface(args)
}

// generateRepositoryInterface is the function that will be executed when the command is called
// It will generate a new repository interface file based on the template
// The template will be parsed with the repositoryData
// Accepts the name of the module as an argument
func generateRepositoryInterface(args []string) {
	templateContent := content.RepositoryInterfaceTemplate
	moduleName := strings.ToLower(args[0])
	fileName := fmt.Sprintf("%s_repository_interface", moduleName)
	modelPath := filepath.Join(projectName, "src/app", moduleName)
	repositoryData := repositoryData{
		Name:      pkg.StringTitle(moduleName),
		LowerName: strings.ToLower(moduleName),
		ModelPath: modelPath,
	}

	if repository.CRUD {
		templateContent = content.RepositoryInterfaceCRUDTemplate
		runMakeModelCommand(nil, args)
	}

	// Generate repository file based on template
	pkg.GenerateFile("repository", fileName, moduleName, templateContent, repositoryData)
}

// generateRepository is the function that will be executed when the command is called
// It will generate a new repository file based on the template
// The template will be parsed with the repositoryData
// Accepts the name of the module as an argument
func generateRepository(args []string) {
	templateContent := content.RepositoryTemplate
	moduleName := strings.ToLower(args[0])
	fileName := fmt.Sprintf("gorm_%s_repository", moduleName)
	modelPath := filepath.Join(projectName, "src/app", moduleName)
	repositoryData := repositoryData{
		Name:      pkg.StringTitle(moduleName),
		LowerName: strings.ToLower(moduleName),
		ModelPath: modelPath,
	}

	if repository.CRUD {
		templateContent = content.RepositoryCRUDTemplate
		runMakeModelCommand(nil, args)
	}

	// Generate repository file based on template
	pkg.GenerateFile("repository", fileName, moduleName, templateContent, repositoryData)
}
