package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/bitvcs/bit/internal/domain"
	sqlcSqlite "github.com/bitvcs/bit/internal/repository/sqlc/sqlite"
)

type User struct {
	queries *sqlcSqlite.Queries
}

func NewUserRepository(db *sql.DB) *User {
	return &User{queries: sqlcSqlite.New(db)}
}

func (a *User) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := a.queries.UserGetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewErrorRecordNotFound()
		}
		return nil, domain.NewErrorDatabase(err.Error())
	}
	return toDomainUser(user), nil
}

func (a *User) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	user, err := a.queries.UserGetById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewErrorRecordNotFound()
		}
		return nil, domain.NewErrorDatabase(err.Error())
	}
	return toDomainUser(user), nil
}

func toDomainUser(row sqlcSqlite.User) *domain.User {
	user := domain.User{
		ID:           row.ID,
		Name:         row.Name,
		Email:        row.Email,
		Password:     row.Password,
		PhotoUrl:     row.PhotoUrl.String,
		IsSuperAdmin: row.IsSuperAdmin,
		IsAdmin:      row.IsAdmin,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
		Deleted:      row.Deleted,
		DeletedAt:    nullTimePtr(row.DeletedAt),
	}
	return &user
}
