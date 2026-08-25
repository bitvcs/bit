package usecase

import (
	"context"

	"github.com/nipalab/nipa/internal/domain"
)

type branchRepository interface {
	ListBranches(ctx context.Context, projectID int64, limit int, updatedAfter *int64, lastID int64) ([]*domain.Branch, error)
}

type Branch struct {
	branchRepo branchRepository
}

func NewBranch(branchRepo branchRepository) *Branch {
	return &Branch{
		branchRepo: branchRepo,
	}
}

func (b *Branch) ListBranches(ctx context.Context, projectID int64, limit int, updatedAfter *int64, lastID int64) ([]*domain.Branch, error) {
	return b.branchRepo.ListBranches(ctx, projectID, limit, updatedAfter, lastID)
}
