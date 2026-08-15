package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/apinprastya/bit/internal/domain"
	"github.com/apinprastya/bit/internal/usecase"
)

// fakeProjectRepository is an in-memory stub implementing usecase.ProjectRepository.
type fakeProjectRepository struct {
	createFn func(ctx context.Context, project domain.Project) (domain.Project, error)
	getFn    func(ctx context.Context, id int64) (domain.Project, error)
	listFn   func(ctx context.Context) ([]domain.Project, error)
	deleteFn func(ctx context.Context, id int64) error
}

func (f *fakeProjectRepository) Create(ctx context.Context, project domain.Project) (domain.Project, error) {
	return f.createFn(ctx, project)
}

func (f *fakeProjectRepository) Get(ctx context.Context, id int64) (domain.Project, error) {
	return f.getFn(ctx, id)
}

func (f *fakeProjectRepository) List(ctx context.Context) ([]domain.Project, error) {
	return f.listFn(ctx)
}

func (f *fakeProjectRepository) Delete(ctx context.Context, id int64) error {
	return f.deleteFn(ctx, id)
}

func TestProjectUsecase_Create(t *testing.T) {
	ctx := context.Background()
	want := domain.Project{ID: 1, Name: "project-a", Description: "a"}

	repo := &fakeProjectRepository{
		createFn: func(_ context.Context, project domain.Project) (domain.Project, error) {
			require.Equal(t, "project-a", project.Name)
			return want, nil
		},
	}

	got, err := usecase.NewProjectUsecase(repo).Create(ctx, domain.Project{Name: "project-a", Description: "a"})
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestProjectUsecase_Create_Error(t *testing.T) {
	ctx := context.Background()
	wantErr := errors.New("create failed")

	repo := &fakeProjectRepository{
		createFn: func(_ context.Context, _ domain.Project) (domain.Project, error) {
			return domain.Project{}, wantErr
		},
	}

	_, err := usecase.NewProjectUsecase(repo).Create(ctx, domain.Project{})
	require.ErrorIs(t, err, wantErr)
}

func TestProjectUsecase_Get(t *testing.T) {
	ctx := context.Background()
	want := domain.Project{ID: 42, Name: "project-a"}

	repo := &fakeProjectRepository{
		getFn: func(_ context.Context, id int64) (domain.Project, error) {
			require.Equal(t, int64(42), id)
			return want, nil
		},
	}

	got, err := usecase.NewProjectUsecase(repo).Get(ctx, 42)
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestProjectUsecase_Get_Error(t *testing.T) {
	ctx := context.Background()
	wantErr := errors.New("not found")

	repo := &fakeProjectRepository{
		getFn: func(_ context.Context, _ int64) (domain.Project, error) {
			return domain.Project{}, wantErr
		},
	}

	_, err := usecase.NewProjectUsecase(repo).Get(ctx, 42)
	require.ErrorIs(t, err, wantErr)
}

func TestProjectUsecase_List(t *testing.T) {
	ctx := context.Background()
	want := []domain.Project{{ID: 1}, {ID: 2}}

	repo := &fakeProjectRepository{
		listFn: func(_ context.Context) ([]domain.Project, error) {
			return want, nil
		},
	}

	got, err := usecase.NewProjectUsecase(repo).List(ctx)
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestProjectUsecase_List_Error(t *testing.T) {
	ctx := context.Background()
	wantErr := errors.New("list failed")

	repo := &fakeProjectRepository{
		listFn: func(_ context.Context) ([]domain.Project, error) {
			return nil, wantErr
		},
	}

	_, err := usecase.NewProjectUsecase(repo).List(ctx)
	require.ErrorIs(t, err, wantErr)
}

func TestProjectUsecase_Delete(t *testing.T) {
	ctx := context.Background()
	called := false

	repo := &fakeProjectRepository{
		deleteFn: func(_ context.Context, id int64) error {
			called = true
			require.Equal(t, int64(7), id)
			return nil
		},
	}

	err := usecase.NewProjectUsecase(repo).Delete(ctx, 7)
	require.NoError(t, err)
	require.True(t, called)
}

func TestProjectUsecase_Delete_Error(t *testing.T) {
	ctx := context.Background()
	wantErr := errors.New("delete failed")

	repo := &fakeProjectRepository{
		deleteFn: func(_ context.Context, _ int64) error {
			return wantErr
		},
	}

	err := usecase.NewProjectUsecase(repo).Delete(ctx, 7)
	require.ErrorIs(t, err, wantErr)
}
