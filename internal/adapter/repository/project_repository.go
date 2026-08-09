// Package repository adapts sqlc-generated queries to the usecase repository ports.
package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/apinprastya/bit/internal/domain"
	"github.com/apinprastya/bit/internal/infrastructure/database/sqlc"
)

// ProjectRepository implements usecase.ProjectRepository backed by sqlite.
type ProjectRepository struct {
	queries *sqlc.Queries
}

func NewProjectRepository(db *sql.DB) *ProjectRepository {
	return &ProjectRepository{queries: sqlc.New(db)}
}

func (r *ProjectRepository) Create(ctx context.Context, project domain.Project) (domain.Project, error) {
	row, err := r.queries.CreateProject(ctx, sqlc.CreateProjectParams{
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

func (r *ProjectRepository) List(ctx context.Context) ([]domain.Project, error) {
	rows, err := r.queries.ListProjects(ctx)
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

func toDomainProject(row sqlc.Project) domain.Project {
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
