package sqlite_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	sqliterepo "github.com/apinprastya/bit/internal/adapter/repository/sqlite"
	"github.com/apinprastya/bit/internal/domain"
	"github.com/apinprastya/bit/internal/infrastructure/database"
)

func TestProjectRepositorySQLite(t *testing.T) {
	db, err := database.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	defer func() { require.NoError(t, db.Close()) }()

	require.NoError(t, database.MigrateUp(db, "sqlite3"))

	repo := sqliterepo.NewProjectRepository(db)
	ctx := context.Background()

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
