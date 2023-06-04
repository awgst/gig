package migrate

import (
	"context"
	"database/sql"
	"fmt"
	"os"
)

// Migrate is the struct that contains the database connection
type Migrate struct {
	DB *sql.DB
}

// NewMigrate creates a new instance of Migrate
func NewMigrate(db *sql.DB) *Migrate {
	return &Migrate{
		DB: db,
	}
}

// Transaction
func (m *Migrate) transaction(ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := m.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	err = fn(tx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	// Commit transaction
	return tx.Commit()
}

// Up migrates the database to the latest version
func (m *Migrate) Up(ctx context.Context, paths map[int]string) error {
	return m.transaction(ctx, func(tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS gig_migrations ( version int(11) NOT NULL )`); err != nil {
			return err
		}

		// Get latest migration version
		rows, err := tx.QueryContext(ctx, `SELECT version FROM gig_migrations order by version desc limit 1`)
		if err != nil {
			return err
		}

		var latestVersion int
		for rows.Next() {
			err = rows.Scan(&latestVersion)
			if err != nil {
				return err
			}
		}

		for k, path := range paths {
			if k <= latestVersion {
				continue
			}

			// Read file
			c, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			// Execute sql
			sql := string(c)
			if _, err := tx.ExecContext(ctx, sql); err != nil {
				return err
			}

			// Insert migration version
			if _, err := tx.ExecContext(ctx, `INSERT INTO gig_migrations (version) VALUES (?)`, k); err != nil {
				return err
			}

			fmt.Printf("Migrating: %s\n", path)
		}

		return nil
	})
}

// Down migrates the database to the previous version
func (m *Migrate) Down(ctx context.Context, paths map[int]string) error {
	return m.transaction(ctx, func(tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS gig_migrations ( version int(11) NOT NULL )`); err != nil {
			return err
		}

		// Get latest migration version
		rows, err := tx.QueryContext(ctx, `SELECT version FROM gig_migrations order by version desc limit 1`)
		if err != nil {
			return err
		}

		var latestVersion int
		for rows.Next() {
			err = rows.Scan(&latestVersion)
			if err != nil {
				return err
			}
		}

		for k, path := range paths {
			if k > latestVersion {
				continue
			}
			// Read file
			c, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			// Execute sql
			sql := string(c)
			if _, err := tx.ExecContext(ctx, sql); err != nil {
				return err
			}

			// Insert migration version
			if _, err := tx.ExecContext(ctx, `DELETE FROM gig_migrations WHERE version = ?`, k); err != nil {
				return err
			}

			fmt.Printf("Rolling Back: %s\n", path)
		}

		return nil
	})
}

// Reset migrates the database to the initial version
func (m *Migrate) Reset() error {
	m.setupTable()
	_, err := m.DB.Exec(`UPDATE gig_migrations SET status = 'pending'`)
	if err != nil {
		return err
	}
	return nil
}

// Make inserts a new migration to the database
func (m *Migrate) Make(version int) error {
	m.setupTable()
	_, err := m.DB.Exec(`INSERT INTO gig_migrations (version, status) VALUES (?, ?)`, version, "applied")
	if err != nil {
		return err
	}
	return nil
}

// Create tables gig_migrations
func (m *Migrate) setupTable() error {
	_, err := m.DB.Exec(`CREATE TABLE IF NOT EXISTS gig_migrations (
		version int(11) NOT NULL
	)`)
	if err != nil {
		return err
	}

	return nil
}

// GetLatestVersion returns the latest version of the migration
func (m *Migrate) GetLatestVersion() int {
	var version int
	err := m.DB.QueryRow(`SELECT version FROM gig_migrations ORDER BY version DESC LIMIT 1`).Scan(&version)
	if err != nil {
		return 0
	}
	return version
}
