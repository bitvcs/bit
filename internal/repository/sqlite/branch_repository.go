package sqlite

import (
	"context"
	"database/sql"
	"time"

	"github.com/nipalab/nipa/internal/domain"
	sqlcSqlite "github.com/nipalab/nipa/internal/repository/sqlc/sqlite"
	"github.com/nipalab/nipa/internal/snow"
	"gopkg.in/typ.v4/slices"
)

type BranchRepository struct {
	queries *sqlcSqlite.Queries
}

func NewBranchRepository(db *sql.DB) *BranchRepository {
	return &BranchRepository{queries: sqlcSqlite.New(db)}
}

func (b *BranchRepository) ListBranches(ctx context.Context, projectID snow.ID, limit int, updatedAfter *time.Time, lastID snow.ID) ([]*domain.Branch, error) {
	rows, err := b.queries.BranchList(ctx, sqlcSqlite.BranchListParams{
		ProjectID:     projectID.Int64(),
		Limit:         int64(limit),
		LastUpdatedAt: timePtrToNullTime(updatedAfter),
		LastID:        lastID.Int64(),
	})
	if err != nil {
		return nil, handleError(err)
	}

	return slices.Map(rows, func(b sqlcSqlite.Branch) *domain.Branch {
		return branchToDomain(b)
	}), nil
}

func (b *BranchRepository) GetByProjectIDAndID(ctx context.Context, projectID snow.ID, branchID snow.ID) (*domain.Branch, error) {
	row, err := b.queries.BranchGet(ctx, sqlcSqlite.BranchGetParams{
		ProjectID: projectID.Int64(),
		ID:        branchID.Int64(),
	})
	if err != nil {
		return nil, handleError(err)
	}
	return branchToDomain(row), nil
}

func branchToDomain(b sqlcSqlite.Branch) *domain.Branch {
	var commitID *snow.ID
	if b.CommitID.Valid {
		id := snow.ID(b.CommitID.Int64)
		commitID = &id
	}
	return &domain.Branch{
		ID:        snow.ID(b.ID),
		ProjectID: snow.ID(b.ProjectID),
		Name:      b.Name,
		Protected: b.Protected,
		CommitID:  commitID,
		UpdatedAt: b.UpdatedAt,
		CreatedAt: b.CreatedAt,
		Deleted:   b.Deleted,
		DeletedAt: nullTimePtr(b.DeletedAt),
	}
}
