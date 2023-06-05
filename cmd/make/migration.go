// Package make implements the command to generate a new file
package make

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	content "github.com/awgst/gig/template/migration"

	"github.com/awgst/gig/pkg"
	"github.com/spf13/cobra"
)

// migration is the data that will be parsed in the template
type migrationData struct {
	TableName string
}

// MigrationOptions is the options that will be parsed in the command
type MigrationOptions struct {
	Path string
}

var migration MigrationOptions

// MigrationCommand is the command to generate a new migration
var MigrationCommand = &cobra.Command{
	Use:   "make:migration <name>",
	Short: "Make a new migration",
	Long:  "Make a new migration",
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
	Run:     runMakeMigrationCommand,
	Example: `gig make:migration create_users_table`,
}

func init() {
	migration = MigrationOptions{}
	flags := MigrationCommand.Flags()
	flags.StringVar(&migration.Path, "path", "", "Specify the path that will be used to migrate")
}

// runMakeMigrationCommand is the function that will be executed when the command is called
// It will generate a new migration file based on the template
// The template will be parsed with the migrationData
// Accepts the name of the module as an argument
func runMakeMigrationCommand(cmd *cobra.Command, args []string) {
	GenerateMigration(migration, args)
}

// GenerateMigration is function that will be called from another package
// Accept MigrationOptions and []string as an argument
func GenerateMigration(migrationOpt MigrationOptions, args []string) {
	var templateUpContent string
	var templateDownContent string
	var tableName string
	fileName := strings.ToLower(args[0])
	if strings.Contains(fileName, "create") {
		tableName = strings.Replace(fileName, "create_", "", -1)
		tableName = strings.Replace(tableName, "_table", "", -1)
		templateUpContent = content.CreateMigrationSQLTemplate
		templateDownContent = content.DropMigrationSQLTemplate
	}
	// Generate migration file based on template
	migrationData := migrationData{
		TableName: tableName,
	}

	// Path
	path := filepath.Join("database", "migrations")
	if migrationOpt.Path != "" {
		path = migrationOpt.Path
	}

	// Get latest migration version from path
	latestVersion := 0
	files, _ := filepath.Glob(filepath.Join(path, "*.up.sql"))
	filesLength := len(files)
	if filesLength > 0 {
		file := files[filesLength-1]
		regex := regexp.MustCompile(`^(\d+)_`)
		version := regex.FindStringSubmatch(filepath.Base(file))
		if len(version) > 0 {
			latestVersion, _ = strconv.Atoi(version[1])
		}
	}

	latestVersion = latestVersion + 1

	pkg.GenerateFile(
		"migration",
		fmt.Sprintf("%06d_%s.up.sql", latestVersion, fileName),
		"",
		templateUpContent,
		migrationData,
		path,
	)
	pkg.GenerateFile(
		"migration",
		fmt.Sprintf("%06d_%s.down.sql", latestVersion, fileName),
		"",
		templateDownContent,
		migrationData,
		path,
	)
}
