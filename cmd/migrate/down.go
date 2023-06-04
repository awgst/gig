// Package make implements the command to migrate files with suffix `.down.sql` and `.down.sql`
package migrate

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/awgst/gig/pkg/database"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

type migrateDownOptions struct {
	Path       string
	Connection string
	From       int
	To         int
}

var down migrateDownOptions

// DownCommand is the sub command down of migrate
var downCommand = &cobra.Command{
	Use:     "down",
	Short:   "Migrate down to the latest version",
	Long:    "Migrate down to the latest version",
	Run:     runMigrateDownCommand,
	Example: `gig migrate down`,
}

func init() {
	godotenv.Load(".env")
	flags := downCommand.Flags()
	flags.StringVar(&down.Path, "path", "", "Specify the path that will be used to migrate")
	flags.IntVar(&down.From, "from", 0, "Specify the version where the migration will start")
	flags.IntVar(&down.To, "to", 0, "Specify the version where the migration will end")
	flags.StringVar(&down.Connection, "connection", "", "Specify the connection that will be used to migrate")
}

// runMigrateDownCommand is the function that will be executed when the command is called
func runMigrateDownCommand(cmd *cobra.Command, args []string) {
	// Path
	path := filepath.Join("database", "migrations")
	if down.Path != "" {
		path = down.Path
	}

	// Get driver and dsn
	driver, dbDsn := getDriverAndDsn()
	if down.Connection != "" {
		dbDsn = down.Connection
	}
	m := NewMigrate(database.ConnectSql(driver, dbDsn))
	defer m.DB.Close()

	// Get map of all migrations file in the path
	paths := getPathForDown(m, path)

	// Down migrations
	if err := m.Down(context.Background(), paths); err != nil {
		log.Fatal(err)
	}

	fmt.Println("All migrations have been applied")
}

// getPathForDown returns the paths of the migrations
func getPathForDown(m *Migrate, path string) map[int]string {
	fromFile := 0
	toFile := 0
	paths := map[int]string{}

	fromFile = m.GetLatestVersion()

	files, _ := filepath.Glob(filepath.Join(path, "*.down.sql"))
	filesLength := len(files)
	if down.From != 0 {
		fromFile = down.From
	}

	if fromFile != 0 {
		toFile = fromFile - 1
	}

	if down.To != 0 {
		toFile = down.To
	}
	for i := filesLength - 1; i >= 0; i-- {
		file := files[i]
		regex := regexp.MustCompile(`^(\d+)_`)
		version := regex.FindStringSubmatch(filepath.Base(file))
		if len(version) > 0 {
			v, _ := strconv.Atoi(version[1])
			if v <= fromFile && v > toFile {
				paths[v] = file
			}
		}
	}

	return paths
}
