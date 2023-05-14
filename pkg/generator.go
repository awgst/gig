package pkg

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/briandowns/spinner"
)

func GenerateFile(templateType string, fileName string, moduleName string, templateContent string, templateData any) {
	// Create template
	tmpl := template.Must(template.New(templateType).Parse(templateContent))

	folderPath := filepath.Join("src/app", moduleName, templateType)
	filePath := filepath.Join(folderPath, fmt.Sprintf("%s.go", fileName))

	// Create the folder if it doesn't exist
	err := os.MkdirAll(folderPath, 0755)
	if err != nil {
		log.Fatal("Failed to create folder:", err)
	}

	// Check if the file already exists
	if _, err := os.Stat(filePath); err == nil {
		return
	}

	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
	s.Suffix = fmt.Sprintf(" Creating %s\n", filePath)
	s.Start()
	// Create a new file to write the generated code
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal("Failed to create file:", err)
	}
	defer file.Close()

	// Execute the template and write the result to the file
	err = tmpl.Execute(file, templateData)
	if err != nil {
		log.Fatal("Failed to generate Go file:", err)
	}
}
