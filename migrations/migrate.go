package migrations

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed data/*.sql
var fs embed.FS

func Migrate(db *sql.DB) error {
	databaseDriver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		return fmt.Errorf("cannot create sqlite instance: %w", err)
	}

	sourceDriver, err := iofs.New(fs, "data")
	if err != nil {
		return fmt.Errorf("cannot create iofs from embeded data: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", sourceDriver, "sqlite", databaseDriver)
	if err != nil {
		return fmt.Errorf("cannot create migrate: %w", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("cannotrun migrations: %w", err)
	}

	return nil
}
