package sqlite

import (
	"context"
	"database/sql"
	"time"

	"github.com/bitvcs/bit/internal/domain"
	sqlcSqlite "github.com/bitvcs/bit/internal/repository/sqlc/sqlite"
)

type ProjectRepository struct {
	queries *sqlcSqlite.Queries
}

func NewProjectRepository(db *sql.DB) *ProjectRepository {
	return &ProjectRepository{queries: sqlcSqlite.New(db)}
}

func (r *ProjectRepository) Create(ctx context.Context, project domain.Project) (domain.Project, error) {
	row, err := r.queries.CreateProject(ctx, sqlcSqlite.CreateProjectParams{
		Name:        project.Name,
		Description: project.Description,
	})
	if err != nil {
		return domain.Project{}, err
	}
	return toDomainProject(row), nil
}

func (r *ProjectRepository) Get(ctx context.Context, id int64) (domain.Project, error) {
	row, err := r.queries.GetProject(ctx, id)
	if err != nil {
		return domain.Project{}, err
	}
	return toDomainProject(row), nil
}

func (r *ProjectRepository) ListByOrgID(ctx context.Context, orgID int64) ([]domain.Project, error) {
	rows, err := r.queries.ListProjectsByOrgId(ctx, orgID)
	if err != nil {
		return nil, err
	}
	projects := make([]domain.Project, 0, len(rows))
	for _, row := range rows {
		projects = append(projects, toDomainProject(row))
	}
	return projects, nil
}

func (r *ProjectRepository) Delete(ctx context.Context, id int64) error {
	return r.queries.DeleteProject(ctx, id)
}

func toDomainProject(row sqlcSqlite.Project) domain.Project {
	return domain.Project{
		ID:          row.ID,
		Name:        row.Name,
		Description: row.Description,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
		DeletedAt:   nullTimePtr(row.DeletedAt),
	}
}

func nullTimePtr(t sql.NullTime) *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}
