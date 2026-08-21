package sqlite

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAuthRepositorySQLite_SaveAndGetAndDeleteRefreshToken(t *testing.T) {
	ctx := context.Background()
	db, q := newSQLiteTestDB(t)
	repo := NewAuthRepository(db)

	userID := seedUser(t, q, "alice", "alice@example.com", sql.NullString{})
	expiresAt := time.Now().Add(60 * 24 * time.Hour)

	require.NoError(t, repo.SaveRefreshToken(ctx, userID, "token-1", expiresAt))

	got, err := repo.GetAndDeleteRefreshToken(ctx, "token-1")
	require.NoError(t, err)
	require.NotZero(t, got.ID)
	require.Equal(t, userID, got.UserID)
	require.Equal(t, "token-1", got.Token)
	require.WithinDuration(t, expiresAt, got.ExpiresAt, time.Second)
	require.False(t, got.CreatedAt.IsZero())
}

func TestAuthRepositorySQLite_RefreshTokenIsSingleUse(t *testing.T) {
	ctx := context.Background()
	db, q := newSQLiteTestDB(t)
	repo := NewAuthRepository(db)

	userID := seedUser(t, q, "bob", "bob@example.com", sql.NullString{})
	require.NoError(t, repo.SaveRefreshToken(ctx, userID, "token-2", time.Now().Add(time.Hour)))

	_, err := repo.GetAndDeleteRefreshToken(ctx, "token-2")
	require.NoError(t, err)

	_, err = repo.GetAndDeleteRefreshToken(ctx, "token-2")
	requireRecordNotFound(t, err)
}

func TestAuthRepositorySQLite_GetAndDeleteUnknownToken(t *testing.T) {
	ctx := context.Background()
	db, _ := newSQLiteTestDB(t)
	repo := NewAuthRepository(db)

	_, err := repo.GetAndDeleteRefreshToken(ctx, "unknown-token")
	requireRecordNotFound(t, err)
}
