package usecase_test

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"

	"github.com/bitvcs/bit/internal/domain"
	"github.com/bitvcs/bit/internal/usecase"
)

// fakeUserRepository is an in-memory stub implementing usecase.userRepository.
type fakeUserRepository struct {
	getByEmailFn func(ctx context.Context, email string) (*domain.User, error)
	getByIDFn    func(ctx context.Context, id int64) (*domain.User, error)
}

func (f *fakeUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return f.getByEmailFn(ctx, email)
}

func (f *fakeUserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	return f.getByIDFn(ctx, id)
}

// fakeAuthRepository is an in-memory stub implementing usecase.authRepository.
type fakeAuthRepository struct {
	saveRefreshTokenFn         func(ctx context.Context, userID int64, refreshToken string, expiresAt int64) error
	getAndDeleteRefreshTokenFn func(ctx context.Context, refreshToken string) (domain.RefreshToken, error)
}

func (f *fakeAuthRepository) SaveRefreshToken(ctx context.Context, userID int64, refreshToken string, expiresAt int64) error {
	return f.saveRefreshTokenFn(ctx, userID, refreshToken, expiresAt)
}

func (f *fakeAuthRepository) GetAndDeleteRefreshToken(ctx context.Context, refreshToken string) (domain.RefreshToken, error) {
	return f.getAndDeleteRefreshTokenFn(ctx, refreshToken)
}

func TestAuth_LoginWithEmailPassword(t *testing.T) {
	ctx := context.Background()
	user := &domain.User{ID: 42, Email: "user@example.com"}

	userRepo := &fakeUserRepository{
		getByEmailFn: func(_ context.Context, email string) (*domain.User, error) {
			require.Equal(t, "user@example.com", email)
			return user, nil
		},
	}
	authRepo := &fakeAuthRepository{
		saveRefreshTokenFn: func(_ context.Context, userID int64, refreshToken string, expiresAt int64) error {
			require.Equal(t, int64(42), userID)
			require.NotEmpty(t, refreshToken)
			require.InDelta(t, time.Now().Add(60*24*time.Hour).Unix(), expiresAt, 5)
			return nil
		},
	}

	got, err := usecase.NewAuth("secret", userRepo, authRepo).LoginWithEmailPassword(ctx, "user@example.com", "password")
	require.NoError(t, err)

	require.Equal(t, "Bearer", got.TokenType)
	require.NotEmpty(t, got.RefreshToken)
	require.NotEmpty(t, got.AccessToken)

	claims := parseAccessToken(t, got.AccessToken, "secret")
	require.Equal(t, strconv.FormatInt(42, 10), claims.Subject)
	require.InDelta(t, 30*time.Minute.Seconds(), claims.ExpiresAt.Sub(claims.IssuedAt.Time).Seconds(), 2)
}

func TestAuth_LoginWithEmailPassword_UserNotFound(t *testing.T) {
	ctx := context.Background()
	wantErr := errors.New("not found")

	userRepo := &fakeUserRepository{
		getByEmailFn: func(_ context.Context, _ string) (*domain.User, error) {
			return nil, wantErr
		},
	}
	authRepo := &fakeAuthRepository{}

	_, err := usecase.NewAuth("secret", userRepo, authRepo).LoginWithEmailPassword(ctx, "user@example.com", "password")
	require.ErrorIs(t, err, wantErr)
}

func TestAuth_LoginWithEmailPassword_SaveRefreshTokenError(t *testing.T) {
	ctx := context.Background()
	wantErr := errors.New("save failed")

	userRepo := &fakeUserRepository{
		getByEmailFn: func(_ context.Context, _ string) (*domain.User, error) {
			return &domain.User{ID: 1}, nil
		},
	}
	authRepo := &fakeAuthRepository{
		saveRefreshTokenFn: func(_ context.Context, _ int64, _ string, _ int64) error {
			return wantErr
		},
	}

	_, err := usecase.NewAuth("secret", userRepo, authRepo).LoginWithEmailPassword(ctx, "user@example.com", "password")
	require.ErrorIs(t, err, wantErr)
}

func TestAuth_LoginWithRefreshToken(t *testing.T) {
	ctx := context.Background()
	stored := domain.RefreshToken{UserID: 7, Token: "old-refresh-token", ExpiresAt: time.Now().Add(time.Hour)}

	authRepo := &fakeAuthRepository{
		getAndDeleteRefreshTokenFn: func(_ context.Context, refreshToken string) (domain.RefreshToken, error) {
			require.Equal(t, "old-refresh-token", refreshToken)
			return stored, nil
		},
		saveRefreshTokenFn: func(_ context.Context, userID int64, refreshToken string, expiresAt int64) error {
			require.Equal(t, int64(7), userID)
			require.NotEqual(t, "old-refresh-token", refreshToken)
			require.InDelta(t, time.Now().Add(60*24*time.Hour).Unix(), expiresAt, 5)
			return nil
		},
	}
	userRepo := &fakeUserRepository{
		getByIDFn: func(_ context.Context, id int64) (*domain.User, error) {
			require.Equal(t, int64(7), id)
			return &domain.User{ID: 7}, nil
		},
	}

	got, err := usecase.NewAuth("secret", userRepo, authRepo).LoginWithRefreshToken(ctx, "old-refresh-token")
	require.NoError(t, err)

	require.Equal(t, "Bearer", got.TokenType)
	require.NotEmpty(t, got.RefreshToken)
	require.NotEmpty(t, got.AccessToken)

	claims := parseAccessToken(t, got.AccessToken, "secret")
	require.Equal(t, strconv.FormatInt(7, 10), claims.Subject)
}

func TestAuth_LoginWithRefreshToken_Expired(t *testing.T) {
	ctx := context.Background()

	authRepo := &fakeAuthRepository{
		getAndDeleteRefreshTokenFn: func(_ context.Context, _ string) (domain.RefreshToken, error) {
			return domain.RefreshToken{UserID: 7, ExpiresAt: time.Now().Add(-time.Hour)}, nil
		},
	}
	userRepo := &fakeUserRepository{}

	_, err := usecase.NewAuth("secret", userRepo, authRepo).LoginWithRefreshToken(ctx, "expired-token")
	require.Error(t, err)

	var domErr *domain.Error
	require.ErrorAs(t, err, &domErr)
	require.Equal(t, 400, domErr.Code)
	require.Equal(t, "refresh token expired", domErr.Message)
}

func TestAuth_LoginWithRefreshToken_GetError(t *testing.T) {
	ctx := context.Background()
	wantErr := errors.New("token not found")

	authRepo := &fakeAuthRepository{
		getAndDeleteRefreshTokenFn: func(_ context.Context, _ string) (domain.RefreshToken, error) {
			return domain.RefreshToken{}, wantErr
		},
	}
	userRepo := &fakeUserRepository{}

	_, err := usecase.NewAuth("secret", userRepo, authRepo).LoginWithRefreshToken(ctx, "unknown-token")
	require.ErrorIs(t, err, wantErr)
}

func TestAuth_LoginWithRefreshToken_UserNotFound(t *testing.T) {
	ctx := context.Background()
	wantErr := errors.New("not found")

	authRepo := &fakeAuthRepository{
		getAndDeleteRefreshTokenFn: func(_ context.Context, _ string) (domain.RefreshToken, error) {
			return domain.RefreshToken{UserID: 99, ExpiresAt: time.Now().Add(time.Hour)}, nil
		},
	}
	userRepo := &fakeUserRepository{
		getByIDFn: func(_ context.Context, _ int64) (*domain.User, error) {
			return nil, wantErr
		},
	}

	_, err := usecase.NewAuth("secret", userRepo, authRepo).LoginWithRefreshToken(ctx, "valid-token")
	require.ErrorIs(t, err, wantErr)
}

func parseAccessToken(t *testing.T, token string, secret string) jwt.RegisteredClaims {
	t.Helper()

	claims := jwt.RegisteredClaims{}
	parsed, err := jwt.ParseWithClaims(token, &claims, func(_ *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	require.NoError(t, err)
	require.True(t, parsed.Valid)
	return claims
}
