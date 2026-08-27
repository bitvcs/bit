package usecase

import (
	"context"
	"time"

	"github.com/nipalab/nipa/internal/domain"
	"github.com/nipalab/nipa/internal/snow"
)

type permissionUsecase interface {
	HasProjectAccess(ctx context.Context, projectID snow.ID, permission domain.Permission) bool
}

type branchRepository interface {
	ListBranches(ctx context.Context, projectID snow.ID, limit int, updatedAfter *time.Time, lastID snow.ID) ([]*domain.Branch, error)
	GetByProjectIDAndID(ctx context.Context, projectID snow.ID, branchID snow.ID) (*domain.Branch, error)
	GetDefaultBranch(ctx context.Context, projectID snow.ID) (*domain.Branch, error)
}

type Branch struct {
	permUc     permissionUsecase
	branchRepo branchRepository
}

func NewBranch(permUc permissionUsecase, branchRepo branchRepository) *Branch {
	return &Branch{
		permUc:     permUc,
		branchRepo: branchRepo,
	}
}

func (b *Branch) ListBranches(ctx context.Context, projectID snow.ID, limit int, updatedAfter *time.Time, lastID snow.ID) ([]*domain.Branch, error) {
	if !b.permUc.HasProjectAccess(ctx, projectID, domain.PermissionRead) {
		return nil, domain.NewErrorNoPermission()
	}
	return b.branchRepo.ListBranches(ctx, projectID, limit, updatedAfter, lastID)
}

func (b *Branch) GetByProjectIDAndID(ctx context.Context, projectID snow.ID, branchID snow.ID) (*domain.Branch, error) {
	if !b.permUc.HasProjectAccess(ctx, projectID, domain.PermissionRead) {
		return nil, domain.NewErrorNoPermission()
	}
	return b.branchRepo.GetByProjectIDAndID(ctx, projectID, branchID)
}

func (b *Branch) GetDefault(ctx context.Context, projectID snow.ID) (*domain.Branch, error) {
	if !b.permUc.HasProjectAccess(ctx, projectID, domain.PermissionRead) {
		return nil, domain.NewErrorNoPermission()
	}
	return b.branchRepo.GetDefaultBranch(ctx, projectID)
}
