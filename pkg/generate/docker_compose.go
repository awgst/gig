// Package generate implements of generate files for command needs
package generate

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	dockercompose "github.com/awgst/gig/template/docker-compose"
)

// Generate docker compose file
// Accepts projectName and database as string
// Returns error
func GenerateDockerComposeFile(projectName, database string) error {
	filename := "docker-compose.yml"

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
		`version: '3'
services:
    backend:
        build: ./
        container_name: 'gig_backend'
        ports:
            - '${APP_PORT}:${APP_PORT}'
        extra_hosts:
            - "host.docker.internal:host-gateway"
        volumes:
            - ./:/app/
            - ./go.mod:/go/src/app/go.mod
        networks:
            - gig
        depends_on:
            - %s
    %s
networks:
    gig:
        driver: bridge
volumes:
    gig-%s:
        driver: local`,
		database,
		getDatabaseForDockerCompose(database),
		database,
	)

	type tmplVars struct {
		Version   string
		CamelName string
	}

	vars := tmplVars{
		Version:   "1",
		CamelName: filename,
	}

	tmpl := template.Must(template.New("docker-compose").Parse(tmplString))
	if err := tmpl.Execute(f, vars); err != nil {
		return err
	}

	return nil
}

// Get database for docker compose file
func getDatabaseForDockerCompose(database string) string {
	databases := map[string]string{
		"mysql":      dockercompose.MySqlDBTemplate,
		"postgresql": dockercompose.PostgreSqlDBTemplate,
	}

	return databases[database]
}
