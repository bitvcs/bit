package postgres_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/bitvcs/bit/internal/adapter/repository/postgres"
	"github.com/bitvcs/bit/internal/domain"
	"github.com/bitvcs/bit/internal/infrastructure/database"
)

var postgresDB *sql.DB

func TestProjectRepositoryPostgres(t *testing.T) {
	ctx := context.Background()
	container, err := tcpostgres.Run(ctx,
		"postgres:16-alpine",
		tcpostgres.WithDatabase("bit_test"),
		tcpostgres.WithUsername("bit"),
		tcpostgres.WithPassword("bit"),
		tcpostgres.BasicWaitStrategies(),
	)
	require.NoError(t, err)
	defer func() { require.NoError(t, container.Terminate(ctx)) }()

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	db, err := database.Open("postgres", connStr)
	require.NoError(t, err)
	defer func() { require.NoError(t, db.Close()) }()

	require.NoError(t, database.MigrateUp(db, "postgres"))
	postgresDB = db

	repo := postgres.NewProjectRepository(db)

	_, err = db.Exec("TRUNCATE TABLE projects RESTART IDENTITY CASCADE")
	require.NoError(t, err)

	created, err := repo.Create(ctx, domain.Project{
		Name:        "project-a",
		Description: "first project",
	})
	require.NoError(t, err)
	require.NotZero(t, created.ID)
	require.Equal(t, "project-a", created.Name)
	require.Equal(t, "first project", created.Description)

	got, err := repo.Get(ctx, created.ID)
	require.NoError(t, err)
	require.Equal(t, created.ID, got.ID)
	require.Equal(t, created.Name, got.Name)
	require.Equal(t, created.Description, got.Description)

	_, err = repo.Get(ctx, 999999)
	require.Error(t, err)

	_, err = repo.Create(ctx, domain.Project{Name: "project-b", Description: "second"})
	require.NoError(t, err)

	projects, err := repo.List(ctx)
	require.NoError(t, err)
	require.Len(t, projects, 2)

	created2, err := repo.Create(ctx, domain.Project{Name: "project-c", Description: "third"})
	require.NoError(t, err)
	require.NoError(t, repo.Delete(ctx, created2.ID))

	_, err = repo.Get(ctx, created2.ID)
	require.Error(t, err)
}
