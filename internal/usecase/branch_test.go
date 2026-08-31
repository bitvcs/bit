package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/nipalab/nipa/internal/domain"
	"github.com/nipalab/nipa/internal/snow"
)

func TestNewBranch(t *testing.T) {
	ctrl := gomock.NewController(t)
	perm := NewMockpermissionUsecase(ctrl)
	repo := NewMockbranchRepository(ctrl)

	uc := NewBranch(perm, repo)
	require.NotNil(t, uc)
}

func TestBranch_ListBranches_NoPermission(t *testing.T) {
	ctrl := gomock.NewController(t)
	perm := NewMockpermissionUsecase(ctrl)
	repo := NewMockbranchRepository(ctrl)

	perm.EXPECT().
		HasProjectAccess(gomock.Any(), snow.ID(1), domain.PermissionRead).
		Return(false)

	uc := NewBranch(perm, repo)
	_, err := uc.ListBranches(context.Background(), snow.ID(1), 10, nil, 0)
	require.Error(t, err)

	var domErr *domain.Error
	require.ErrorAs(t, err, &domErr)
	require.Equal(t, 403, domErr.Code)
}

func TestBranch_ListBranches_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	perm := NewMockpermissionUsecase(ctrl)
	repo := NewMockbranchRepository(ctrl)

	want := []*domain.Branch{
		{ID: 1, ProjectID: 1, Name: "main"},
		{ID: 2, ProjectID: 1, Name: "develop"},
	}

	perm.EXPECT().
		HasProjectAccess(gomock.Any(), snow.ID(1), domain.PermissionRead).
		Return(true)

	after := time.Now().Add(-time.Hour)
	lastID := snow.ID(0)
	repo.EXPECT().
		ListBranches(gomock.Any(), snow.ID(1), 10, &after, lastID).
		Return(want, nil)

	uc := NewBranch(perm, repo)
	got, err := uc.ListBranches(context.Background(), snow.ID(1), 10, &after, lastID)
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestBranch_ListBranches_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	perm := NewMockpermissionUsecase(ctrl)
	repo := NewMockbranchRepository(ctrl)

	wantErr := errors.New("db down")

	perm.EXPECT().
		HasProjectAccess(gomock.Any(), snow.ID(1), domain.PermissionRead).
		Return(true)

	repo.EXPECT().
		ListBranches(gomock.Any(), snow.ID(1), 10, gomock.Nil(), snow.ID(0)).
		Return(nil, wantErr)

	uc := NewBranch(perm, repo)
	_, err := uc.ListBranches(context.Background(), snow.ID(1), 10, nil, 0)
	require.ErrorIs(t, err, wantErr)
}

func TestBranch_ListBranches_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	perm := NewMockpermissionUsecase(ctrl)
	repo := NewMockbranchRepository(ctrl)

	perm.EXPECT().
		HasProjectAccess(gomock.Any(), snow.ID(1), domain.PermissionRead).
		Return(true)

	repo.EXPECT().
		ListBranches(gomock.Any(), snow.ID(1), 10, gomock.Nil(), snow.ID(0)).
		Return([]*domain.Branch{}, nil)

	uc := NewBranch(perm, repo)
	got, err := uc.ListBranches(context.Background(), snow.ID(1), 10, nil, 0)
	require.NoError(t, err)
	require.Empty(t, got)
}

func TestBranch_GetByProjectIDAndID_NoPermission(t *testing.T) {
	ctrl := gomock.NewController(t)
	perm := NewMockpermissionUsecase(ctrl)
	repo := NewMockbranchRepository(ctrl)

	perm.EXPECT().
		HasProjectAccess(gomock.Any(), snow.ID(1), domain.PermissionRead).
		Return(false)

	uc := NewBranch(perm, repo)
	_, err := uc.GetByProjectIDAndID(context.Background(), snow.ID(1), snow.ID(2))
	require.Error(t, err)

	var domErr *domain.Error
	require.ErrorAs(t, err, &domErr)
	require.Equal(t, 403, domErr.Code)
}

func TestBranch_GetByProjectIDAndID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	perm := NewMockpermissionUsecase(ctrl)
	repo := NewMockbranchRepository(ctrl)

	want := &domain.Branch{ID: 2, ProjectID: 1, Name: "main"}

	perm.EXPECT().
		HasProjectAccess(gomock.Any(), snow.ID(1), domain.PermissionRead).
		Return(true)

	repo.EXPECT().
		GetByProjectIDAndID(gomock.Any(), snow.ID(1), snow.ID(2)).
		Return(want, nil)

	uc := NewBranch(perm, repo)
	got, err := uc.GetByProjectIDAndID(context.Background(), snow.ID(1), snow.ID(2))
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestBranch_GetByProjectIDAndID_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	perm := NewMockpermissionUsecase(ctrl)
	repo := NewMockbranchRepository(ctrl)

	wantErr := errors.New("db down")

	perm.EXPECT().
		HasProjectAccess(gomock.Any(), snow.ID(1), domain.PermissionRead).
		Return(true)

	repo.EXPECT().
		GetByProjectIDAndID(gomock.Any(), snow.ID(1), snow.ID(2)).
		Return(nil, wantErr)

	uc := NewBranch(perm, repo)
	_, err := uc.GetByProjectIDAndID(context.Background(), snow.ID(1), snow.ID(2))
	require.ErrorIs(t, err, wantErr)
}

func TestBranch_GetByProjectIDAndID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	perm := NewMockpermissionUsecase(ctrl)
	repo := NewMockbranchRepository(ctrl)

	perm.EXPECT().
		HasProjectAccess(gomock.Any(), snow.ID(1), domain.PermissionRead).
		Return(true)

	repo.EXPECT().
		GetByProjectIDAndID(gomock.Any(), snow.ID(1), snow.ID(2)).
		Return(nil, domain.NewErrorRecordNotFound())

	uc := NewBranch(perm, repo)
	_, err := uc.GetByProjectIDAndID(context.Background(), snow.ID(1), snow.ID(2))
	require.Error(t, err)

	var domErr *domain.Error
	require.ErrorAs(t, err, &domErr)
	require.Equal(t, 404, domErr.Code)
}

func TestBranch_GetDefault_NoPermission(t *testing.T) {
	ctrl := gomock.NewController(t)
	perm := NewMockpermissionUsecase(ctrl)
	repo := NewMockbranchRepository(ctrl)

	perm.EXPECT().
		HasProjectAccess(gomock.Any(), snow.ID(1), domain.PermissionRead).
		Return(false)

	uc := NewBranch(perm, repo)
	_, err := uc.GetDefault(context.Background(), snow.ID(1))
	require.Error(t, err)

	var domErr *domain.Error
	require.ErrorAs(t, err, &domErr)
	require.Equal(t, 403, domErr.Code)
}

func TestBranch_GetDefault_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	perm := NewMockpermissionUsecase(ctrl)
	repo := NewMockbranchRepository(ctrl)

	want := &domain.Branch{ID: 1, ProjectID: 1, Name: "main", IsDefault: true}

	perm.EXPECT().
		HasProjectAccess(gomock.Any(), snow.ID(1), domain.PermissionRead).
		Return(true)

	repo.EXPECT().
		GetDefaultBranch(gomock.Any(), snow.ID(1)).
		Return(want, nil)

	uc := NewBranch(perm, repo)
	got, err := uc.GetDefault(context.Background(), snow.ID(1))
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestBranch_GetDefault_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	perm := NewMockpermissionUsecase(ctrl)
	repo := NewMockbranchRepository(ctrl)

	wantErr := errors.New("db down")

	perm.EXPECT().
		HasProjectAccess(gomock.Any(), snow.ID(1), domain.PermissionRead).
		Return(true)

	repo.EXPECT().
		GetDefaultBranch(gomock.Any(), snow.ID(1)).
		Return(nil, wantErr)

	uc := NewBranch(perm, repo)
	_, err := uc.GetDefault(context.Background(), snow.ID(1))
	require.ErrorIs(t, err, wantErr)
}
