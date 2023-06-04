// Package pkg implements list function and variable that can be used by other packages
package pkg

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/briandowns/spinner"
)

// GenerateFile generates a file from a template
// Accepts templateType, fileName, moduleName, templateContent, templateData
func GenerateFile(templateType string, fileName string, moduleName string, templateContent string, templateData any, path ...string) {
	// Create template
	tmpl := template.Must(template.New(templateType).Parse(templateContent))

	folderPath := filepath.Join("src/app", moduleName, templateType)
	if len(path) > 0 {
		folderPath = path[0]
	}
	if !strings.Contains(fileName, ".sql") {
		fileName = fmt.Sprintf("%s.go", fileName)
	}
	filePath := filepath.Join(folderPath, fileName)

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
	s.Stop()
}
