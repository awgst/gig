// Package generate implements of generate files for command needs
package generate

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
)

// Generate json file
// Accepts projectName, database and httpFramework as string
// Returns error
func GenerateJsonFile(projectName, database, httpFramework string) error {
	filename := "gig.json"

	fullPath := filepath.Join(projectName, filename)
	if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
		return err
	}

	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer f.Close()

	var tmplString = fmt.Sprintf(
		`{
	"version": "1",
	"name": "%s",
	"database": "%s",
	"http_framework": "%s"
}`,
		projectName,
		database,
		httpFramework,
	)

	type tmplVars struct {
		Version   string
		CamelName string
	}

	vars := tmplVars{
		Version:   "1",
		CamelName: filename,
	}

	tmpl := template.Must(template.New("gig").Parse(tmplString))
	if err := tmpl.Execute(f, vars); err != nil {
		return err
	}

	return nil
}
