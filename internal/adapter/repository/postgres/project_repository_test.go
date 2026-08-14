package postgres_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/suite"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/apinprastya/bit/internal/adapter/repository/postgres"
	"github.com/apinprastya/bit/internal/domain"
	"github.com/apinprastya/bit/internal/infrastructure/database"
)

// ProjectRepositorySuite shares a single postgres container across all its tests.
type ProjectRepositorySuite struct {
	suite.Suite

	container *tcpostgres.PostgresContainer
	db        *sql.DB
	repo      *postgres.ProjectRepository
}

func TestProjectRepositorySuite(t *testing.T) {
	suite.Run(t, new(ProjectRepositorySuite))
}

func (s *ProjectRepositorySuite) SetupSuite() {
	ctx := context.Background()

	container, err := tcpostgres.Run(ctx,
		"postgres:16-alpine",
		tcpostgres.WithDatabase("bit_test"),
		tcpostgres.WithUsername("bit"),
		tcpostgres.WithPassword("bit"),
		tcpostgres.BasicWaitStrategies(),
	)
	s.Require().NoError(err)
	s.container = container

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	s.Require().NoError(err)

	db, err := database.Open("postgres", connStr)
	s.Require().NoError(err)
	s.Require().NoError(database.MigrateUp(db, "postgres"))
	s.db = db

	s.repo = postgres.NewProjectRepository(db)
}

func (s *ProjectRepositorySuite) TearDownSuite() {
	s.Require().NoError(s.db.Close())
	s.Require().NoError(s.container.Terminate(context.Background()))
}

// SetupTest truncates the projects table so each test starts from a clean state.
func (s *ProjectRepositorySuite) SetupTest() {
	_, err := s.db.Exec("TRUNCATE TABLE projects RESTART IDENTITY CASCADE")
	s.Require().NoError(err)
}

func (s *ProjectRepositorySuite) TestCreateAndGet() {
	ctx := context.Background()

	created, err := s.repo.Create(ctx, domain.Project{
		Name:        "project-a",
		Description: "first project",
	})
	s.Require().NoError(err)
	s.Require().NotZero(created.ID)
	s.Equal("project-a", created.Name)
	s.Equal("first project", created.Description)

	got, err := s.repo.Get(ctx, created.ID)
	s.Require().NoError(err)
	s.Equal(created.ID, got.ID)
	s.Equal(created.Name, got.Name)
	s.Equal(created.Description, got.Description)
}

func (s *ProjectRepositorySuite) TestGet_NotFound() {
	ctx := context.Background()

	_, err := s.repo.Get(ctx, 999999)
	s.Error(err)
}

func (s *ProjectRepositorySuite) TestList() {
	ctx := context.Background()

	_, err := s.repo.Create(ctx, domain.Project{Name: "project-a", Description: "a"})
	s.Require().NoError(err)
	_, err = s.repo.Create(ctx, domain.Project{Name: "project-b", Description: "b"})
	s.Require().NoError(err)

	projects, err := s.repo.List(ctx)
	s.Require().NoError(err)
	s.Len(projects, 2)
}

func (s *ProjectRepositorySuite) TestDelete() {
	ctx := context.Background()

	created, err := s.repo.Create(ctx, domain.Project{Name: "project-a", Description: "a"})
	s.Require().NoError(err)

	s.Require().NoError(s.repo.Delete(ctx, created.ID))

	_, err = s.repo.Get(ctx, created.ID)
	s.Error(err)
}
