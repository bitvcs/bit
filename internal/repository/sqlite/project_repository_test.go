package sqlite

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	database "github.com/nipalab/nipa/db"
	"github.com/nipalab/nipa/internal/domain"
	"github.com/nipalab/nipa/internal/snow"
)

func TestProjectRepositorySQLite(t *testing.T) {
	db, err := database.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	defer func() { require.NoError(t, db.Close()) }()

	require.NoError(t, database.MigrateUp(db, "sqlite3"))

	repo := NewProjectRepository(db)
	ctx := context.Background()

	snowNode, err := snow.NewNode(1)
	require.NoError(t, err)
	created, err := repo.Create(ctx, domain.Project{
		ID:          snowNode.Generate(),
		OrgID:       1,
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

	_, err = repo.Create(ctx, domain.Project{ID: snowNode.Generate(), OrgID: 1, Name: "project-b", Description: "second"})
	require.NoError(t, err)

	projects, err := repo.ListByOrgID(ctx, 1)
	require.NoError(t, err)
	require.Len(t, projects, 3)

	created2, err := repo.Create(ctx, domain.Project{ID: snowNode.Generate(), OrgID: 1, Name: "project-c", Description: "third"})
	require.NoError(t, err)
	require.NoError(t, repo.Delete(ctx, created2.ID))

	_, err = repo.Get(ctx, created2.ID)
	require.Error(t, err)
}
