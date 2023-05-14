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

type repositoryData struct {
	Name      string
	LowerName string
	ModelPath string
}

type repositoryOptions struct {
	CRUD   bool
	Module string
}

var repository repositoryOptions

var projectName string

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

func runMakeRepositoryCommand(cmd *cobra.Command, args []string) {
	generateRepository(args)
	generateRepositoryInterface(args)
}

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
