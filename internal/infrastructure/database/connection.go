// Package database provides the sqlite connection and migration runner.
package database

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "modernc.org/sqlite"
)

// driverNames maps a logical dialect name to its registered database/sql driver name.
var driverNames = map[string]string{
	"sqlite3":  "sqlite",
	"postgres": "pgx",
}

// Open opens a database connection for the given dialect ("sqlite3" or "postgres").
func Open(dialect, dataSourceName string) (*sql.DB, error) {
	driverName, ok := driverNames[dialect]
	if !ok {
		return nil, fmt.Errorf("unsupported dialect %q", dialect)
	}

	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}
