package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/awgst/gig/pkg"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

var goVersion = pkg.LatestGoVersion
var projectName string

var createCommand = &cobra.Command{
	Use:     "create <name>/[go-version]",
	Aliases: []string{},
	Short:   "Create a new project",
	Long:    "Create a new project with specific version of go.\nSupported Go version is 1.18, 1.19, 1.20",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a name argument")
		}

		argsSplitted := strings.Split(args[0], "/")
		if len(argsSplitted) > 1 && argsSplitted[1] != "latest" {
			goVersion = argsSplitted[1]
		}

		if !slices.Contains(pkg.SupportedGoVersions, goVersion) {
			return fmt.Errorf("go v%s is unsupported", goVersion)
		}

		projectName = argsSplitted[0]

		return nil
	},
	ValidArgs:  []string{"name"},
	ArgAliases: []string{"name"},
	Example:    "gig create go-project/1.19",
	Run:        runCreateCommand,
}

func init() {
}

func runCreateCommand(cmd *cobra.Command, args []string) {
	err := survey.Ask(
		[]*survey.Question{
			{
				Name: "http_framework",
				Prompt: &survey.Select{
					Message: "Choose a http framework:",
					Options: []string{
						"chi",
						"echo",
						"fiber",
						"gin",
						"mux",
					},
					Default:  "chi",
					PageSize: 10,
				},
				Validate: survey.Required,
			},
			{
				Name: "database",
				Prompt: &survey.Select{
					Message: "Choose database",
					Options: []string{
						"mysql",
						"postgresql",
						"sqlite",
						"sqlserver",
					},
					Default:  "mysql",
					PageSize: 10,
				},
				Validate: survey.Required,
			},
		},
		&createCommandAnswer,
	)
	if err != nil {
		log.Fatal(err)
	}

	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
	s.Suffix = " Creating project...\n"
	s.Start()
	err = pkg.GitClone(projectName, createCommandAnswer.HttpFramework)
	if err != nil {
		log.Fatal(err)
	}

	pkg.Replace(
		projectName,
		fmt.Sprintf("gig-%s-template", createCommandAnswer.HttpFramework),
		projectName,
	)

	goModPath := filepath.Join(projectName, "go.mod")
	pkg.Replace(
		goModPath,
		"go 1.19",
		fmt.Sprintf("go %s", goVersion),
	)

	dockerFilePath := filepath.Join(projectName, "Dockerfile")
	pkg.Replace(
		dockerFilePath,
		"FROM golang:latest",
		fmt.Sprintf("FROM golang:%s-alpine", goVersion),
	)

	err = generateDockerComposeFile(projectName, createCommandAnswer.Database)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(`
📗 Project created successfully. Happy coding!

❗ Please run go mod download && go mod tidy after project created successfully
❔ More informations --> https://github.com/awgst/gig`)
}

// Generate docker compose file
func generateDockerComposeFile(projectName, database string) error {
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
		"mysql": `mysql:
        image: 'mysql:latest'
        container_name: 'gig_mysql'
        ports:
            - '${FORWARD_DB_PORT:-3306}:3306'
        environment:
            MYSQL_ROOT_PASSWORD: '${DB_PASSWORD}'
            MYSQL_ROOT_HOST: "%"
            MYSQL_DATABASE: '${DB_DATABASE}'
            MYSQL_ALLOW_EMPTY_PASSWORD: 'yes'
        volumes:
            - 'gig-mysql:/var/lib/mysql'
            - './database/create-database.sh:/docker-entrypoint-initdb.d/10-create-testing-database.sh'
        networks:
            - gig
        healthcheck:
            test: ["CMD", "mysqladmin", "ping", "-p${DB_PASSWORD}"]
            retries: 3
            timeout: 5s`,
	}

	return databases[database]
}
