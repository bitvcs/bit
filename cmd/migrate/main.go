// Command migrate applies or reverts schema migrations against sqlite or postgres.
package main

import (
	"flag"
	"log"

	"github.com/apinprastya/bit/internal/infrastructure/database"
)

func main() {
	dialect := flag.String("dialect", "sqlite3", "database dialect: sqlite3 or postgres")
	dsn := flag.String("dsn", "bit.db", "data source name (sqlite file path or postgres connection string)")
	action := flag.String("action", "up", "migration action: up or down")
	flag.Parse()

	db, err := database.Open(*dialect, *dsn)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("close database: %v", err)
		}
	}()

	switch *action {
	case "up":
		err = database.MigrateUp(db, *dialect)
	case "down":
		err = database.MigrateDown(db, *dialect)
	default:
		log.Fatalf("unknown action %q (expected up or down)", *action)
	}
	if err != nil {
		log.Fatalf("migrate %s: %v", *action, err)
	}
	log.Printf("migrate %s completed", *action)
}
