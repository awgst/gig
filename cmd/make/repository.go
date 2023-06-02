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
	Name          string
	LowerName     string
	CamelCaseName string
	ModelName     string
	ModelPath     string
	TableName     string
}

// repositoryOptions is the options that will be parsed in the command
type RepositoryOptions struct {
	CRUD   bool
	Module string
	ORM    bool
}

// repository is the options that will be parsed in the command
var repository RepositoryOptions

// projectName is the name of the project
var projectName string

// RepositoryCommand is the command to generate a new repository
var RepositoryCommand = &cobra.Command{
	Use:   "make:repository <name>",
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
	flags.BoolVar(&repository.ORM, "orm", false, "Create a repository with ORM methods using GORM")
	projectName, _ = pkg.ReadJsonString("name")
}

// runMakeRepositoryCommand is the function that will be executed when the command is called
// It will generate a new repository file based on the template
// The template will be parsed with the repositoryData
// Accepts the name of the module as an argument
func runMakeRepositoryCommand(cmd *cobra.Command, args []string) {
	GenerateRepository(repository, args)
}

// GenerateRepository is the function that wrap generateRepository and generateRepositoryInterface
// Accept RepositoryOptions as an argument
func GenerateRepository(repositoryOptions RepositoryOptions, args []string) {
	repository = repositoryOptions
	if repository.CRUD {
		GenerateModel(ModelOptions{
			Module: repository.Module,
			Plain:  false,
		}, args)
	}
	generateRepository(args)
	generateRepositoryInterface(args)
}

// generateRepositoryInterface is the function that will be executed when the command is called
// It will generate a new repository interface file based on the template
// The template will be parsed with the repositoryData
// Accepts the name of the module as an argument
func generateRepositoryInterface(args []string) {
	templateContent := content.RepositoryInterfaceTemplate
	name := strings.ToLower(args[0])
	moduleName := name
	if repository.Module != "" {
		moduleName = repository.Module
	}
	fileName := fmt.Sprintf("%s_repository_interface", name)
	modelPath := filepath.Join(projectName, "src/app", moduleName)
	repositoryData := repositoryData{
		Name:          pkg.SnakeToPascal(name),
		LowerName:     strings.ToLower(name),
		CamelCaseName: pkg.SnakeToCamel(name),
		ModelName:     pkg.SnakeToPascal(name),
		ModelPath:     modelPath,
	}

	if repository.CRUD {
		templateContent = content.RepositoryInterfaceCRUDTemplate
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
	name := strings.ToLower(args[0])
	moduleName := name
	types := "sql"
	if repository.Module != "" {
		moduleName = repository.Module
	}

	if repository.CRUD {
		templateContent = content.RepositoryCRUDTemplate
	}

	useOrm, _ := pkg.ReadJsonBool("use_orm")
	if repository.ORM || useOrm {
		templateContent = content.RepositoryCRUDGormTemplate
		types = "gorm"
	}

	fileName := fmt.Sprintf("%s_%s_repository", types, name)
	modelPath := filepath.Join(projectName, "src/app", moduleName)
	repositoryData := repositoryData{
		Name:          pkg.SnakeToPascal(name),
		LowerName:     strings.ToLower(name),
		CamelCaseName: pkg.SnakeToCamel(name),
		ModelName:     pkg.SnakeToPascal(name),
		ModelPath:     modelPath,
		TableName:     pkg.PluralizeSnakeCase(name),
	}

	// Generate repository file based on template
	pkg.GenerateFile("repository", fileName, moduleName, templateContent, repositoryData)
}
