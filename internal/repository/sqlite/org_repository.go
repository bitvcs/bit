package sqlite

import (
	"context"

	"github.com/nipalab/nipa/internal/domain"
	sqlcSqlite "github.com/nipalab/nipa/internal/repository/sqlc/sqlite"
	"github.com/nipalab/nipa/internal/snow"
)

type OrgRepository struct {
	queries *sqlcSqlite.Queries
}

func NewOrgRepository(queries *sqlcSqlite.Queries) *OrgRepository {
	return &OrgRepository{
		queries: queries,
	}
}

func (r *OrgRepository) GetBySlug(ctx context.Context, slug string) (*domain.Organization, error) {
	org, err := r.queries.GetOrganizationBySlug(ctx, slug)
	if err != nil {
		return nil, handleError(err)
	}
	return &domain.Organization{
		ID:        snow.ID(org.ID),
		Slug:      org.Slug,
		Name:      org.Name,
		CreatedAt: org.CreatedAt,
		UpdatedAt: org.UpdatedAt,
	}, nil
}
