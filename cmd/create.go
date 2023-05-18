// Package cmd implements the list of commands that can be executed
package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/awgst/gig/pkg"
	"github.com/awgst/gig/pkg/generate"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

// goVersion is the version of go that will be used
// by default it's the latest version from pkg.LatestGoVersion
var goVersion = pkg.LatestGoVersion
var projectName string

// createCommand is the command to create a new project
var createCommand = &cobra.Command{
	Use:     "create <name>",
	Aliases: []string{},
	Short:   "Create a new project",
	Long:    "Create a new project",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a name argument")
		}

		projectName = args[0]

		return nil
	},
	ValidArgs:  []string{"name"},
	ArgAliases: []string{"name"},
	Example:    "gig create go-project",
	Run:        runCreateCommand,
}

func init() {
	flags := createCommand.Flags()
	flags.StringVarP(&goVersion, "version", "v", pkg.LatestGoVersion, "Specify version of go. Currently supported Go version is 1.18, 1.19, 1.20")
}

// runCreateCommand is the function that will be executed when the command is called
func runCreateCommand(cmd *cobra.Command, args []string) {
	// Ask the user to choose the http framework and the database
	err := survey.Ask(
		pkg.CreateSurveyQuestion,
		&createCommandAnswer,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Check if the go version is supported
	if !slices.Contains(pkg.SupportedGoVersions, goVersion) {
		log.Fatalf("go v%s is unsupported", goVersion)
	}

	// Create project by cloning the template
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
	s.Suffix = " Creating project...\n"
	s.Start()
	err = pkg.GitClone(projectName, createCommandAnswer.HttpFramework)
	if err != nil {
		log.Fatal(err)
	}

	// Replace the name of the project in the template
	pkg.Replace(
		projectName,
		fmt.Sprintf("gig-%s-template", createCommandAnswer.HttpFramework),
		projectName,
	)

	// Replace the go version in the template
	goModPath := filepath.Join(projectName, "go.mod")
	pkg.Replace(
		goModPath,
		"go 1.19",
		fmt.Sprintf("go %s", goVersion),
	)

	// Replace the go version in the template
	dockerFilePath := filepath.Join(projectName, "Dockerfile")
	pkg.Replace(
		dockerFilePath,
		"FROM golang:latest",
		fmt.Sprintf("FROM golang:%s-alpine", goVersion),
	)

	// Generate the docker-compose file
	err = generate.GenerateDockerComposeFile(projectName, createCommandAnswer.Database)
	if err != nil {
		log.Fatal(err)
	}

	// Generate the json file
	err = generate.GenerateJsonFile(projectName, createCommandAnswer.Database, createCommandAnswer.HttpFramework)
	if err != nil {
		log.Fatal(err)
	}

	goModDownload := exec.Command("go", "mod", "download")
	goModDownload.Dir = projectName
	err = goModDownload.Run()
	if err != nil {
		log.Fatal(err)
	}

	goModTidy := exec.Command("go", "mod", "tidy")
	goModTidy.Dir = projectName
	err = goModTidy.Run()
	if err != nil {
		log.Fatal(err)
	}

	tput := exec.Command("tput", "reset")
	tput.Run()

	fmt.Println(`
📗 Project created successfully
❔ More informations --> https://github.com/awgst/gig`)
}
