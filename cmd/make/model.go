package make

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"regexp"

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
	templateContent := content.ModelTemplate
	moduleName := args[0]
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
	// Create template
	tmpl := template.Must(template.New("model").Parse(templateContent))

	folderPath := filepath.Join("src/app", moduleName, "model")
	filePath := filepath.Join(folderPath, fmt.Sprintf("%s.go", args[0]))

	// Create the folder if it doesn't exist
	err := os.MkdirAll(folderPath, 0755)
	if err != nil {
		log.Fatal("Failed to create folder:", err)
	}

	// Check if the file already exists
	if _, err := os.Stat(filePath); err == nil {
		log.Fatalf("File already exists in %s", filePath)
	}

	fmt.Printf("Creating %s\n", filePath)
	// Create a new file to write the generated code
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal("Failed to create file:", err)
	}
	defer file.Close()

	// Execute the template and write the result to the file
	err = tmpl.Execute(file, modelData)
	if err != nil {
		log.Fatal("Failed to generate Go file:", err)
	}

	fmt.Println("Done!")
}
