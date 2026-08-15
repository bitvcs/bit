// Package migrations embeds the postgres SQL migration files for use by the migration runner.
package migrations

import "embed"

//go:embed *.sql
var FS embed.FS
