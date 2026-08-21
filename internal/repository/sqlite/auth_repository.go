package sqlite

import (
	"context"
	"database/sql"
	"time"

	"github.com/bitvcs/bit/internal/domain"
	sqlcSqlite "github.com/bitvcs/bit/internal/repository/sqlc/sqlite"
)

type Auth struct {
	queries *sqlcSqlite.Queries
}

func NewAuthRepository(db *sql.DB) *Auth {
	return &Auth{queries: sqlcSqlite.New(db)}
}

func (a *Auth) SaveRefreshToken(ctx context.Context, userID int64, refreshToken string, expiresAt time.Time) error {
	_, err := a.queries.RefreshTokenCreate(ctx, sqlcSqlite.RefreshTokenCreateParams{
		UserID:    userID,
		Token:     refreshToken,
		ExpiresAt: expiresAt,
	})
	return handleError(err)
}

func (a *Auth) GetAndDeleteRefreshToken(ctx context.Context, refreshToken string) (*domain.RefreshToken, error) {
	row, err := a.queries.RefreshTokenDeleteByToken(ctx, refreshToken)
	if err != nil {
		return nil, handleError(err)
	}
	return toDomainRefreshToken(row), nil
}

func toDomainRefreshToken(row sqlcSqlite.RefreshToken) *domain.RefreshToken {
	return &domain.RefreshToken{
		ID:        row.ID,
		UserID:    row.UserID,
		Token:     row.Token,
		ExpiresAt: row.ExpiresAt,
		CreatedAt: row.CreatedAt,
	}
}
