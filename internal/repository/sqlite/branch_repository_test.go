package sqlite

import (
	"context"
	"database/sql"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sqlcSqlite "github.com/nipalab/nipa/internal/repository/sqlc/sqlite"
	"github.com/nipalab/nipa/internal/snow"
)

var testNodeCounter int64

func newTestNode(t *testing.T) snow.Node {
	t.Helper()
	id := atomic.AddInt64(&testNodeCounter, 1) % 256
	node, err := snow.NewNode(id)
	require.NoError(t, err)
	return node
}

func seedProject(t *testing.T, q *sqlcSqlite.Queries, orgID int64, name string) snow.ID {
	t.Helper()

	node := newTestNode(t)
	id := node.Generate()

	_, err := q.CreateProject(context.Background(), sqlcSqlite.CreateProjectParams{
		ID:          id.Int64(),
		OrgID:       orgID,
		Slug:        name,
		Name:        name,
		Description: name,
	})
	require.NoError(t, err)
	return id
}

func seedBranch(t *testing.T, db *sql.DB, projectID snow.ID, name string, commitID sql.NullInt64) snow.ID {
	t.Helper()

	node := newTestNode(t)
	id := node.Generate()

	_, err := db.ExecContext(context.Background(),
		`INSERT INTO branches (id, project_id, name, key, commit_id) VALUES (?, ?, ?, ?, ?)`,
		id.Int64(), projectID.Int64(), name, name, commitID,
	)
	require.NoError(t, err)
	return id
}

func TestBranchRepositorySQLite_GetByProjectIDAndID(t *testing.T) {
	ctx := context.Background()
	db, q := newSQLiteTestDB(t)
	repo := NewBranchRepository(db)

	projectID := seedProject(t, q, 1, "test-project")
	branchID := seedBranch(t, db, projectID, "main", sql.NullInt64{})

	got, err := repo.GetByProjectIDAndID(ctx, projectID, branchID)
	require.NoError(t, err)
	require.Equal(t, branchID, got.ID)
	require.Equal(t, projectID, got.ProjectID)
	require.Equal(t, "main", got.Name)
	require.False(t, got.Protected)
	require.Nil(t, got.CommitID)
	require.False(t, got.UpdatedAt.IsZero())
	require.False(t, got.CreatedAt.IsZero())
	require.False(t, got.Deleted)
	require.Nil(t, got.DeletedAt)
}

func TestBranchRepositorySQLite_GetByProjectIDAndID_WithCommitID(t *testing.T) {
	ctx := context.Background()
	db, q := newSQLiteTestDB(t)
	repo := NewBranchRepository(db)

	projectID := seedProject(t, q, 1, "test-project")

	node := newTestNode(t)
	commitID := node.Generate()

	branchID := seedBranch(t, db, projectID, "feature", sql.NullInt64{Int64: commitID.Int64(), Valid: true})

	got, err := repo.GetByProjectIDAndID(ctx, projectID, branchID)
	require.NoError(t, err)
	require.NotNil(t, got.CommitID)
	require.Equal(t, commitID, *got.CommitID)
}

func TestBranchRepositorySQLite_GetByProjectIDAndID_NotFound(t *testing.T) {
	ctx := context.Background()
	db, q := newSQLiteTestDB(t)
	repo := NewBranchRepository(db)

	projectID := seedProject(t, q, 1, "test-project")

	_, err := repo.GetByProjectIDAndID(ctx, projectID, 999999)
	requireRecordNotFound(t, err)
}

func TestBranchRepositorySQLite_GetByProjectIDAndID_WrongProject(t *testing.T) {
	ctx := context.Background()
	db, q := newSQLiteTestDB(t)
	repo := NewBranchRepository(db)

	projectA := seedProject(t, q, 1, "project-a")
	projectB := seedProject(t, q, 1, "project-b")
	branchID := seedBranch(t, db, projectA, "main", sql.NullInt64{})

	_, err := repo.GetByProjectIDAndID(ctx, projectB, branchID)
	requireRecordNotFound(t, err)
}

func TestBranchRepositorySQLite_ListBranches(t *testing.T) {
	ctx := context.Background()
	db, q := newSQLiteTestDB(t)
	repo := NewBranchRepository(db)

	projectID := seedProject(t, q, 1, "test-project")
	seedBranch(t, db, projectID, "main", sql.NullInt64{})
	seedBranch(t, db, projectID, "develop", sql.NullInt64{})
	seedBranch(t, db, projectID, "feature", sql.NullInt64{})

	branches, err := repo.ListBranches(ctx, projectID, 10, nil, 0)
	require.NoError(t, err)
	require.Len(t, branches, 3)
}

func TestBranchRepositorySQLite_ListBranches_EmptyProject(t *testing.T) {
	ctx := context.Background()
	db, q := newSQLiteTestDB(t)
	repo := NewBranchRepository(db)

	projectID := seedProject(t, q, 1, "empty-project")

	branches, err := repo.ListBranches(ctx, projectID, 10, nil, 0)
	require.NoError(t, err)
	require.Empty(t, branches)
}

func TestBranchRepositorySQLite_ListBranches_Limit(t *testing.T) {
	ctx := context.Background()
	db, q := newSQLiteTestDB(t)
	repo := NewBranchRepository(db)

	projectID := seedProject(t, q, 1, "test-project")
	seedBranch(t, db, projectID, "branch-1", sql.NullInt64{})
	seedBranch(t, db, projectID, "branch-2", sql.NullInt64{})
	seedBranch(t, db, projectID, "branch-3", sql.NullInt64{})

	branches, err := repo.ListBranches(ctx, projectID, 2, nil, 0)
	require.NoError(t, err)
	require.Len(t, branches, 2)
}

func TestBranchRepositorySQLite_ListBranches_IsolatedPerProject(t *testing.T) {
	ctx := context.Background()
	db, q := newSQLiteTestDB(t)
	repo := NewBranchRepository(db)

	projectA := seedProject(t, q, 1, "project-a")
	projectB := seedProject(t, q, 1, "project-b")
	seedBranch(t, db, projectA, "main", sql.NullInt64{})
	seedBranch(t, db, projectA, "develop", sql.NullInt64{})
	seedBranch(t, db, projectB, "main", sql.NullInt64{})

	branchesA, err := repo.ListBranches(ctx, projectA, 10, nil, 0)
	require.NoError(t, err)
	require.Len(t, branchesA, 2)

	branchesB, err := repo.ListBranches(ctx, projectB, 10, nil, 0)
	require.NoError(t, err)
	require.Len(t, branchesB, 1)
	require.Equal(t, "main", branchesB[0].Name)
}

func TestBranchRepositorySQLite_ListBranches_Pagination(t *testing.T) {
	ctx := context.Background()
	db, q := newSQLiteTestDB(t)
	repo := NewBranchRepository(db)

	projectID := seedProject(t, q, 1, "test-project")

	now := time.Now().Truncate(time.Second)
	for i := 0; i < 5; i++ {
		node := newTestNode(t)
		id := node.Generate()
		ts := now.Add(time.Duration(i) * time.Minute)

		_, err := db.ExecContext(ctx,
			`INSERT INTO branches (id, project_id, name, key, updated_at, created_at) VALUES (?, ?, ?, ?, ?, ?)`,
			id.Int64(), projectID.Int64(), "branch-"+string(rune('a'+i)), "branch-"+string(rune('a'+i)), ts, ts,
		)
		require.NoError(t, err)
	}

	firstPage, err := repo.ListBranches(ctx, projectID, 2, nil, 0)
	require.NoError(t, err)
	require.Len(t, firstPage, 2)

	lastBranch := firstPage[len(firstPage)-1]
	secondPage, err := repo.ListBranches(ctx, projectID, 2, &lastBranch.UpdatedAt, lastBranch.ID)
	require.NoError(t, err)
	require.Len(t, secondPage, 2)

	for _, b := range secondPage {
		require.True(t, b.UpdatedAt.Before(lastBranch.UpdatedAt) ||
			(b.UpdatedAt.Equal(lastBranch.UpdatedAt) && b.ID < lastBranch.ID))
	}
}
