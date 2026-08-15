package postgres_test

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/apinprastya/bit/internal/adapter/repository/postgres"
	sqliteadapter "github.com/apinprastya/bit/internal/adapter/repository/sqlite"
	"github.com/apinprastya/bit/internal/domain"
	"github.com/apinprastya/bit/internal/infrastructure/database"
)

var postgresDB *sql.DB

type projectRepository interface {
	Create(ctx context.Context, project domain.Project) (domain.Project, error)
	Get(ctx context.Context, id int64) (domain.Project, error)
	List(ctx context.Context) ([]domain.Project, error)
	Delete(ctx context.Context, id int64) error
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	container, err := tcpostgres.Run(ctx,
		"postgres:16-alpine",
		tcpostgres.WithDatabase("bit_test"),
		tcpostgres.WithUsername("bit"),
		tcpostgres.WithPassword("bit"),
		tcpostgres.BasicWaitStrategies(),
	)
	if err != nil {
		panic(err)
	}

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		panic(err)
	}

	db, err := database.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	if err := database.MigrateUp(db, "postgres"); err != nil {
		panic(err)
	}
	postgresDB = db

	code := m.Run()

	if err := db.Close(); err != nil {
		panic(err)
	}
	if err := container.Terminate(ctx); err != nil {
		panic(err)
	}

	os.Exit(code)
}

func TestProjectRepository(t *testing.T) {
	types := []struct {
		name string
		new  func(t *testing.T) (projectRepository, func())
	}{
		{
			name: "postgres",
			new: func(t *testing.T) (projectRepository, func()) {
				t.Helper()
				_, err := postgresDB.Exec("TRUNCATE TABLE projects RESTART IDENTITY CASCADE")
				require.NoError(t, err)
				return postgres.NewProjectRepository(postgresDB), func() {}
			},
		},
		{
			name: "sqlite",
			new: func(t *testing.T) (projectRepository, func()) {
				t.Helper()
				db, err := database.Open("sqlite3", ":memory:")
				require.NoError(t, err)
				require.NoError(t, database.MigrateUp(db, "sqlite3"))
				return sqliteadapter.NewProjectRepository(db), func() {
					require.NoError(t, db.Close())
				}
			},
		},
	}

	for _, tc := range types {
		t.Run(tc.name, func(t *testing.T) {
			repo, cleanup := tc.new(t)
			defer cleanup()

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
		})
	}
}
