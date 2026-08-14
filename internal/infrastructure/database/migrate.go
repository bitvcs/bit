// Package database provides the sqlite/postgres connections and golang-migrate based migration runner.
package database

import (
	"database/sql"
	"errors"
	"fmt"
	"io/fs"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	pgxmigrate "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	postgresmigrations "github.com/apinprastya/bit/db/migrations/postgres"
	sqlitemigrations "github.com/apinprastya/bit/db/migrations/sqlite"
)

func newMigrator(db *sql.DB, dialect string) (*migrate.Migrate, error) {
	var (
		migrationsFS fs.FS
		dbDriver     database.Driver
		err          error
	)

	switch dialect {
	case "sqlite3":
		migrationsFS = sqlitemigrations.FS
		dbDriver, err = sqlite.WithInstance(db, &sqlite.Config{})
	case "postgres":
		migrationsFS = postgresmigrations.FS
		dbDriver, err = pgxmigrate.WithInstance(db, &pgxmigrate.Config{})
	default:
		return nil, fmt.Errorf("unsupported dialect %q", dialect)
	}
	if err != nil {
		return nil, err
	}

	sourceDriver, err := iofs.New(migrationsFS, ".")
	if err != nil {
		return nil, err
	}

	return migrate.NewWithInstance("iofs", sourceDriver, dialect, dbDriver)
}

// MigrateUp applies every migration that has not been applied yet, in order.
func MigrateUp(db *sql.DB, dialect string) error {
	m, err := newMigrator(db, dialect)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}

// MigrateDown reverts every applied migration, from the most recent to the first.
func MigrateDown(db *sql.DB, dialect string) error {
	m, err := newMigrator(db, dialect)
	if err != nil {
		return err
	}
	if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}
