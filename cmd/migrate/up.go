// Package make implements the command to migrate files with suffix `.up.sql` and `.down.sql`
package migrate

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/awgst/gig/pkg/database"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

type migrateUpOptions struct {
	Path       string
	Connection string
	From       int
	To         int
}

var up migrateUpOptions

// UpCommand is the sub command up of migrate
var upCommand = &cobra.Command{
	Use:     "up",
	Short:   "Migrate up to the latest version",
	Long:    "Migrate up to the latest version",
	Run:     runMigrateUpCommand,
	Example: `gig migrate up`,
}

func init() {
	godotenv.Load(".env")
	flags := upCommand.Flags()
	flags.StringVar(&up.Path, "path", "", "Specify the path that will be used to migrate")
	flags.IntVar(&up.From, "from", 0, "Specify the version where the migration will start")
	flags.IntVar(&up.To, "to", 0, "Specify the version where the migration will end")
	flags.StringVar(&up.Connection, "connection", "", "Specify the connection that will be used to migrate")
}

// runMigrateUpCommand is the function that will be executed when the command is called
func runMigrateUpCommand(cmd *cobra.Command, args []string) {
	// Path
	path := filepath.Join("database", "migrations")
	if up.Path != "" {
		path = up.Path
	}

	// Get driver and dsn
	_, dbDsn := getDriverAndDsn()
	if up.Connection != "" {
		dbDsn = up.Connection
	}

	// Create migrate instance
	m, err := migrate.New(
		fmt.Sprintf("file://%v", path),
		dbDsn,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Migrate up
	if err := m.Up(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("All migrations have been applied")
}

// getDriverAndDsn returns the driver and the dsn
func getDriverAndDsn() (string, string) {
	driver := os.Getenv("DB_DRIVER")
	dbHost := os.Getenv("FORWARD_DB_HOST")
	dbPort := os.Getenv("FORWARD_DB_PORT")
	if dbHost == "" {
		dbHost = os.Getenv("DB_HOST")
	}
	if dbPort == "" {
		dbPort = os.Getenv("DB_PORT")
	}
	dbDsn := database.GetDsn(
		driver,
		dbHost,
		dbPort,
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_DATABASE"),
	)

	return driver, dbDsn
}

// getPaths returns the paths of the migrations
// func getPaths(path string) map[int]string {
// 	toFile := 0
// 	paths := map[int]string{}

// 	files, _ := filepath.Glob(filepath.Join(path, "*.up.sql"))
// 	filesLength := len(files)
// 	lastFile := files[filesLength-1]
// 	regex := regexp.MustCompile(`^(\d+)_`)
// 	version := regex.FindStringSubmatch(filepath.Base(lastFile))
// 	if len(version) > 0 {
// 		toFile, _ = strconv.Atoi(version[1])
// 	}
// 	if up.To != 0 {
// 		toFile = up.To
// 	}
// 	for i := 0; i < filesLength; i++ {
// 		file := files[i]
// 		regex := regexp.MustCompile(`^(\d+)_`)
// 		version := regex.FindStringSubmatch(filepath.Base(file))
// 		if len(version) > 0 {
// 			v, _ := strconv.Atoi(version[1])
// 			if v >= up.From && v <= toFile {
// 				paths[v] = file
// 			}
// 		}
// 	}

// 	return paths
// }
