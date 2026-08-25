package server

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/nipalab/nipa/internal/domain"
	"github.com/nipalab/nipa/internal/grpc/pb"
	"github.com/nipalab/nipa/internal/snow"
	"github.com/nipalab/nipa/internal/usecase"
)

type mockUsecaseContainer struct {
	branch *usecase.Branch
}

func (m *mockUsecaseContainer) Auth() *usecase.Auth   { return nil }
func (m *mockUsecaseContainer) User() *usecase.User   { return nil }
func (m *mockUsecaseContainer) Branch() *usecase.Branch { return m.branch }

func mustBase36(id snow.ID) string {
	return id.Base36()
}

func newTestBranchUc(t *testing.T) (*usecase.Branch, *MockpermissionUsecase, *MockbranchRepository) {
	t.Helper()
	ctrl := gomock.NewController(t)
	perm := NewMockpermissionUsecase(ctrl)
	repo := NewMockbranchRepository(ctrl)
	return usecase.NewBranch(perm, repo), perm, repo
}

func TestNew(t *testing.T) {
	branch, _, _ := newTestBranchUc(t)
	srv := New(&mockUsecaseContainer{branch: branch})
	require.NotNil(t, srv)
}

func TestGetListBranch_Success(t *testing.T) {
	branch, perm, repo := newTestBranchUc(t)
	srv := New(&mockUsecaseContainer{branch: branch})

	projectID := snow.ID(42)
	now := time.Now().Truncate(time.Second)
	branchID := snow.ID(99)
	commitID := snow.ID(7)

	perm.EXPECT().
		HasProjectAccess(gomock.Any(), projectID, domain.PermissionRead).
		Return(true)

	repo.EXPECT().
		ListBranches(gomock.Any(), projectID, 10, gomock.Any(), snow.ID(0)).
		Return([]*domain.Branch{
			{
				ID:        branchID,
				ProjectID: projectID,
				Name:      "main",
				Protected: true,
				CommitID:  &commitID,
				UpdatedAt: now,
				CreatedAt: now,
			},
		}, nil)

	resp, err := srv.GetListBranch(context.Background(), &pb.GetListBranchRequest{
		Context: &pb.ProjectContext{ProjectId: mustBase36(projectID)},
		Limit:   10,
	})
	require.NoError(t, err)
	require.Len(t, resp.Branches, 1)

	b := resp.Branches[0]
	require.Equal(t, branchID.Base36(), b.Id)
	require.Equal(t, "main", b.Name)
	require.True(t, b.Protected)
	require.Equal(t, commitID.Base36(), b.CommitId)
	require.Equal(t, now.Unix(), b.CreatedAt.AsTime().Unix())
	require.Equal(t, now.Unix(), b.UpdatedAt.AsTime().Unix())
}

func TestGetListBranch_Empty(t *testing.T) {
	branch, perm, repo := newTestBranchUc(t)
	srv := New(&mockUsecaseContainer{branch: branch})

	projectID := snow.ID(42)

	perm.EXPECT().
		HasProjectAccess(gomock.Any(), projectID, domain.PermissionRead).
		Return(true)

	repo.EXPECT().
		ListBranches(gomock.Any(), projectID, 10, gomock.Nil(), snow.ID(0)).
		Return([]*domain.Branch{}, nil)

	resp, err := srv.GetListBranch(context.Background(), &pb.GetListBranchRequest{
		Context: &pb.ProjectContext{ProjectId: mustBase36(projectID)},
		Limit:   10,
	})
	require.NoError(t, err)
	require.Empty(t, resp.Branches)
}

func TestGetListBranch_WithPagination(t *testing.T) {
	branch, perm, repo := newTestBranchUc(t)
	srv := New(&mockUsecaseContainer{branch: branch})

	projectID := snow.ID(42)
	lastID := snow.ID(50)
	lastUpdate := time.Now().Truncate(time.Second)

	perm.EXPECT().
		HasProjectAccess(gomock.Any(), projectID, domain.PermissionRead).
		Return(true)

	repo.EXPECT().
		ListBranches(gomock.Any(), projectID, 5, gomock.Any(), lastID).
		Return([]*domain.Branch{}, nil)

	lastIDStr := lastID.Base36()
	resp, err := srv.GetListBranch(context.Background(), &pb.GetListBranchRequest{
		Context:       &pb.ProjectContext{ProjectId: mustBase36(projectID)},
		Limit:         5,
		LastUpdatedAt: timePtrToTimestamp(&lastUpdate),
		LastId:        &lastIDStr,
	})
	require.NoError(t, err)
	require.Empty(t, resp.Branches)
}

func TestGetListBranch_InvalidProjectID(t *testing.T) {
	srv := New(&mockUsecaseContainer{branch: nil})

	_, err := srv.GetListBranch(context.Background(), &pb.GetListBranchRequest{
		Context: &pb.ProjectContext{ProjectId: "!@#"},
	})
	require.Error(t, err)
}

func TestGetListBranch_InvalidLastID(t *testing.T) {
	branch, _, _ := newTestBranchUc(t)
	srv := New(&mockUsecaseContainer{branch: branch})

	invalidID := "!!!invalid!!!"
	_, err := srv.GetListBranch(context.Background(), &pb.GetListBranchRequest{
		Context: &pb.ProjectContext{ProjectId: mustBase36(snow.ID(1))},
		LastId:  &invalidID,
	})
	require.Error(t, err)
}

func TestGetListBranch_UsecaseError(t *testing.T) {
	branch, perm, repo := newTestBranchUc(t)
	srv := New(&mockUsecaseContainer{branch: branch})

	projectID := snow.ID(42)

	perm.EXPECT().
		HasProjectAccess(gomock.Any(), projectID, domain.PermissionRead).
		Return(true)

	repo.EXPECT().
		ListBranches(gomock.Any(), projectID, 10, gomock.Nil(), snow.ID(0)).
		Return(nil, domain.NewErrorNoPermission())

	_, err := srv.GetListBranch(context.Background(), &pb.GetListBranchRequest{
		Context: &pb.ProjectContext{ProjectId: mustBase36(projectID)},
		Limit:   10,
	})
	require.Error(t, err)
}

func TestGetBranch_Success(t *testing.T) {
	branch, perm, repo := newTestBranchUc(t)
	srv := New(&mockUsecaseContainer{branch: branch})

	projectID := snow.ID(42)
	branchID := snow.ID(99)
	now := time.Now().Truncate(time.Second)
	commitID := snow.ID(7)

	perm.EXPECT().
		HasProjectAccess(gomock.Any(), projectID, domain.PermissionRead).
		Return(true)

	repo.EXPECT().
		GetByProjectIDAndID(gomock.Any(), projectID, branchID).
		Return(&domain.Branch{
			ID:        branchID,
			ProjectID: projectID,
			Name:      "develop",
			Protected: false,
			CommitID:  &commitID,
			UpdatedAt: now,
			CreatedAt: now,
		}, nil)

	resp, err := srv.GetBranch(context.Background(), &pb.GetBranchRequest{
		Context:  &pb.ProjectContext{ProjectId: mustBase36(projectID)},
		BranchId: mustBase36(branchID),
	})
	require.NoError(t, err)
	require.NotNil(t, resp.Branch)
	require.Equal(t, branchID.Base36(), resp.Branch.Id)
	require.Equal(t, "develop", resp.Branch.Name)
	require.False(t, resp.Branch.Protected)
	require.Equal(t, commitID.Base36(), resp.Branch.CommitId)
}

func TestGetBranch_NotFound(t *testing.T) {
	branch, perm, repo := newTestBranchUc(t)
	srv := New(&mockUsecaseContainer{branch: branch})

	projectID := snow.ID(42)
	branchID := snow.ID(99)

	perm.EXPECT().
		HasProjectAccess(gomock.Any(), projectID, domain.PermissionRead).
		Return(true)

	repo.EXPECT().
		GetByProjectIDAndID(gomock.Any(), projectID, branchID).
		Return(nil, domain.NewErrorRecordNotFound())

	_, err := srv.GetBranch(context.Background(), &pb.GetBranchRequest{
		Context:  &pb.ProjectContext{ProjectId: mustBase36(projectID)},
		BranchId: mustBase36(branchID),
	})
	require.Error(t, err)
}

func TestGetBranch_InvalidProjectID(t *testing.T) {
	srv := New(&mockUsecaseContainer{branch: nil})

	_, err := srv.GetBranch(context.Background(), &pb.GetBranchRequest{
		Context:  &pb.ProjectContext{ProjectId: "!@#"},
		BranchId: "1",
	})
	require.Error(t, err)
}

func TestGetBranch_InvalidBranchID(t *testing.T) {
	srv := New(&mockUsecaseContainer{branch: nil})

	_, err := srv.GetBranch(context.Background(), &pb.GetBranchRequest{
		Context:  &pb.ProjectContext{ProjectId: mustBase36(snow.ID(1))},
		BranchId: "!@#",
	})
	require.Error(t, err)
}

func TestGetBranch_UsecaseError(t *testing.T) {
	branch, perm, repo := newTestBranchUc(t)
	srv := New(&mockUsecaseContainer{branch: branch})

	projectID := snow.ID(42)
	branchID := snow.ID(99)

	perm.EXPECT().
		HasProjectAccess(gomock.Any(), projectID, domain.PermissionRead).
		Return(true)

	repo.EXPECT().
		GetByProjectIDAndID(gomock.Any(), projectID, branchID).
		Return(nil, errors.New("db down"))

	_, err := srv.GetBranch(context.Background(), &pb.GetBranchRequest{
		Context:  &pb.ProjectContext{ProjectId: mustBase36(projectID)},
		BranchId: mustBase36(branchID),
	})
	require.Error(t, err)
}

func timePtrToTimestamp(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}

