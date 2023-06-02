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

	// Replace database connection
	useOrm := "sql"
	replacer := `SQL: database.ConnectSql(env.Get("DB_DRIVER"), Database()[env.Get("DB_DRIVER")]),`
	replaced := `// SQL: database.ConnectSql(env.Get("DB_DRIVER"), Database()[env.Get("DB_DRIVER")]),`
	if createCommandAnswer.UseOrm {
		replacer = `Gorm: database.ConnectGorm(env.Get("DB_DRIVER"), Database()[env.Get("DB_DRIVER")]),`
		replaced = `// Gorm: database.ConnectGorm(env.Get("DB_DRIVER"), Database()[env.Get("DB_DRIVER")]),`
		useOrm = "orm"
	}
	pkg.Replace(
		filepath.Join(projectName, "cmd", "main.go"),
		replaced,
		replacer,
	)

	replaceDatabase(createCommandAnswer.Database, useOrm)

	// Generate the docker-compose file
	err = generate.GenerateDockerComposeFile(projectName, createCommandAnswer.Database)
	if err != nil {
		log.Fatal(err)
	}

	// Generate the json file
	err = generate.GenerateJsonFile(
		projectName,
		createCommandAnswer.Database,
		createCommandAnswer.HttpFramework,
		createCommandAnswer.UseOrm,
	)
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

	fmt.Println("\n📗 Project created successfully")
	fmt.Println("❔ More informations --> https://github.com/awgst/gig")
	fmt.Println("")
	s.Stop()
}

// Replace database connection
func replaceDatabase(database string, conn string) {
	databasesReplaced := map[string]map[string]string{
		"postgresql": {
			"orm": `// "postgres": postgres.Open(dsn),`,
			"sql": `// _ "github.com/lib/pq"`,
		},
		"mysql": {
			"orm": `// "mysql": mysql.Open(dsn),`,
			"sql": `// _ "github.com/go-sql-driver/mysql"`,
		},
	}
	databaseReplacer := map[string]map[string]string{
		"postgresql": {
			"orm": `"postgres": postgres.Open(dsn),`,
			"sql": `_ "github.com/lib/pq"`,
		},
		"mysql": {
			"orm": `"mysql": mysql.Open(dsn),`,
			"sql": `_ "github.com/go-sql-driver/mysql"`,
		},
	}
	fileName := map[string]string{
		"orm": "gorm.go",
		"sql": "sql.go",
	}

	replaced := databasesReplaced[database][conn]
	replacer := databaseReplacer[database][conn]
	pkg.Replace(
		filepath.Join(projectName, "pkg", "database", fileName[conn]),
		replaced,
		replacer,
	)
}
