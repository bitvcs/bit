package sqlite

import (
	"context"
	"database/sql"
	"regexp"
	"strings"

	"github.com/bitvcs/bit/internal/domain"
	sqlcSqlite "github.com/bitvcs/bit/internal/repository/sqlc/sqlite"
)

var nonAlphaNum = regexp.MustCompile(`[^a-z0-9]+`)

type ProjectRepository struct {
	queries *sqlcSqlite.Queries
}

func NewProjectRepository(db *sql.DB) *ProjectRepository {
	return &ProjectRepository{queries: sqlcSqlite.New(db)}
}

func (r *ProjectRepository) Create(ctx context.Context, project domain.Project) (domain.Project, error) {
	slug := project.Slug
	if slug == "" {
		slug = slugify(project.Name)
	}
	row, err := r.queries.CreateProject(ctx, sqlcSqlite.CreateProjectParams{
		OrgID:       project.OrgID,
		Slug:        slug,
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
		OrgID:       row.OrgID,
		Slug:        row.Slug,
		Name:        row.Name,
		Description: row.Description,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
		DeletedAt:   nullTimePtr(row.DeletedAt),
	}
}

func slugify(s string) string {
	s = strings.ToLower(s)
	s = nonAlphaNum.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}
