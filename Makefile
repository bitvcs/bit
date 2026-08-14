SQLC       ?= go run github.com/sqlc-dev/sqlc/cmd/sqlc@latest
MIGRATE_CMD := go run ./cmd/migrate
DRIVER     ?= sqlite3
DSN        ?= bit.db
MIGRATIONS_DIR := db/migrations

.PHONY: sqlc migrate-up migrate-down migrate-create build test

## Generate Go code from SQL queries via sqlc.
sqlc:
	$(SQLC) generate

## Apply all pending migrations. Override DRIVER=postgres DSN=... for postgres.
migrate-up:
	$(MIGRATE_CMD) -dialect $(DRIVER) -dsn $(DSN) -action up

## Revert all applied migrations. Override DRIVER=postgres DSN=... for postgres.
migrate-down:
	$(MIGRATE_CMD) -dialect $(DRIVER) -dsn $(DSN) -action down

## Create a new pair of up/down migration files, e.g. `make migrate-create name=add_users`.
migrate-create:
	@test -n "$(name)" || (echo "usage: make migrate-create name=<migration_name>" && exit 1)
	@ts=$$(date +%Y%m%d%H%M%S); \
	touch "$(MIGRATIONS_DIR)/$${ts}_$(name).up.sql" "$(MIGRATIONS_DIR)/$${ts}_$(name).down.sql"; \
	echo "created $(MIGRATIONS_DIR)/$${ts}_$(name).up.sql and .down.sql"

build:
	go build -o bin/bitd ./cmd/bitd

## Run all tests (requires Docker for testcontainers-backed repository tests).
test:
	go test ./... -v
