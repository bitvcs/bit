package usecase

import (
	"context"

	"github.com/bitvcs/bit/internal/domain"
)

// ProjectRepository is the port implemented by adapters that persist projects.
type ProjectRepository interface {
	Create(ctx context.Context, project domain.Project) (domain.Project, error)
	Get(ctx context.Context, id int64) (domain.Project, error)
	List(ctx context.Context) ([]domain.Project, error)
	Delete(ctx context.Context, id int64) error
}

// ProjectUsecase implements project-related business logic.
type ProjectUsecase struct {
	repo ProjectRepository
}

func NewProjectUsecase(repo ProjectRepository) *ProjectUsecase {
	return &ProjectUsecase{repo: repo}
}

func (u *ProjectUsecase) Create(ctx context.Context, project domain.Project) (domain.Project, error) {
	return u.repo.Create(ctx, project)
}

func (u *ProjectUsecase) Get(ctx context.Context, id int64) (domain.Project, error) {
	return u.repo.Get(ctx, id)
}

func (u *ProjectUsecase) List(ctx context.Context) ([]domain.Project, error) {
	return u.repo.List(ctx)
}

func (u *ProjectUsecase) Delete(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}
