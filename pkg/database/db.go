package database

import (
	"database/sql"
	"fmt"
	"log"

	// Uncomment the code based on your database driver
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func ConnectSql(driver, dbDsn string) *sql.DB {
	// Open database connection
	db, err := sql.Open(driver, dbDsn)
	if err != nil {
		log.Fatal("Database connection error", err)
	}

	return db
}

func GetDsn(driver, host, port, username, password, db string) string {
	return map[string]string{
		"mysql":    fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, db),
		"postgres": fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s", host, port, username, db, password),
	}[driver]
}
