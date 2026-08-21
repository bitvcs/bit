package usecase

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/bitvcs/bit/internal/domain"
)

func TestAuth_LoginWithEmailPassword(t *testing.T) {
	ctx := context.Background()
	user := &domain.User{ID: 42, Email: "user@example.com"}

	ctrl := gomock.NewController(t)
	userRepo := NewMockuserRepository(ctrl)
	authRepo := NewMockauthRepository(ctrl)

	userRepo.EXPECT().
		GetByEmail(gomock.Any(), "user@example.com").
		Return(user, nil)

	authRepo.EXPECT().
		SaveRefreshToken(gomock.Any(), int64(42), gomock.Not(""), gomock.Any()).
		DoAndReturn(func(_ context.Context, _ int64, _ string, expiresAt int64) error {
			require.InDelta(t, time.Now().Add(60*24*time.Hour).Unix(), expiresAt, 5)
			return nil
		})

	got, err := NewAuth("secret", userRepo, authRepo).LoginWithEmailPassword(ctx, "user@example.com", "password")
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

	ctrl := gomock.NewController(t)
	userRepo := NewMockuserRepository(ctrl)
	authRepo := NewMockauthRepository(ctrl)

	userRepo.EXPECT().
		GetByEmail(gomock.Any(), "user@example.com").
		Return(nil, wantErr)

	_, err := NewAuth("secret", userRepo, authRepo).LoginWithEmailPassword(ctx, "user@example.com", "password")
	require.ErrorIs(t, err, wantErr)
}

func TestAuth_LoginWithEmailPassword_SaveRefreshTokenError(t *testing.T) {
	ctx := context.Background()
	wantErr := errors.New("save failed")

	ctrl := gomock.NewController(t)
	userRepo := NewMockuserRepository(ctrl)
	authRepo := NewMockauthRepository(ctrl)

	userRepo.EXPECT().
		GetByEmail(gomock.Any(), "user@example.com").
		Return(&domain.User{ID: 1}, nil)
	authRepo.EXPECT().
		SaveRefreshToken(gomock.Any(), int64(1), gomock.Any(), gomock.Any()).
		Return(wantErr)

	_, err := NewAuth("secret", userRepo, authRepo).LoginWithEmailPassword(ctx, "user@example.com", "password")
	require.ErrorIs(t, err, wantErr)
}

func TestAuth_LoginWithRefreshToken(t *testing.T) {
	ctx := context.Background()
	stored := domain.RefreshToken{UserID: 7, Token: "old-refresh-token", ExpiresAt: time.Now().Add(time.Hour)}

	ctrl := gomock.NewController(t)
	userRepo := NewMockuserRepository(ctrl)
	authRepo := NewMockauthRepository(ctrl)

	authRepo.EXPECT().
		GetAndDeleteRefreshToken(gomock.Any(), "old-refresh-token").
		Return(stored, nil)

	userRepo.EXPECT().
		GetByID(gomock.Any(), int64(7)).
		Return(&domain.User{ID: 7}, nil)

	authRepo.EXPECT().
		SaveRefreshToken(gomock.Any(), int64(7), gomock.Not("old-refresh-token"), gomock.Any()).
		DoAndReturn(func(_ context.Context, _ int64, _ string, expiresAt int64) error {
			require.InDelta(t, time.Now().Add(60*24*time.Hour).Unix(), expiresAt, 5)
			return nil
		})

	got, err := NewAuth("secret", userRepo, authRepo).LoginWithRefreshToken(ctx, "old-refresh-token")
	require.NoError(t, err)

	require.Equal(t, "Bearer", got.TokenType)
	require.NotEmpty(t, got.RefreshToken)
	require.NotEmpty(t, got.AccessToken)

	claims := parseAccessToken(t, got.AccessToken, "secret")
	require.Equal(t, strconv.FormatInt(7, 10), claims.Subject)
}

func TestAuth_LoginWithRefreshToken_Expired(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	userRepo := NewMockuserRepository(ctrl)
	authRepo := NewMockauthRepository(ctrl)

	authRepo.EXPECT().
		GetAndDeleteRefreshToken(gomock.Any(), "expired-token").
		Return(domain.RefreshToken{UserID: 7, ExpiresAt: time.Now().Add(-time.Hour)}, nil)

	_, err := NewAuth("secret", userRepo, authRepo).LoginWithRefreshToken(ctx, "expired-token")
	require.Error(t, err)

	var domErr *domain.Error
	require.ErrorAs(t, err, &domErr)
	require.Equal(t, 400, domErr.Code)
	require.Equal(t, "refresh token expired", domErr.Message)
}

func TestAuth_LoginWithRefreshToken_GetError(t *testing.T) {
	ctx := context.Background()
	wantErr := errors.New("token not found")

	ctrl := gomock.NewController(t)
	userRepo := NewMockuserRepository(ctrl)
	authRepo := NewMockauthRepository(ctrl)

	authRepo.EXPECT().
		GetAndDeleteRefreshToken(gomock.Any(), "unknown-token").
		Return(domain.RefreshToken{}, wantErr)

	_, err := NewAuth("secret", userRepo, authRepo).LoginWithRefreshToken(ctx, "unknown-token")
	require.ErrorIs(t, err, wantErr)
}

func TestAuth_LoginWithRefreshToken_UserNotFound(t *testing.T) {
	ctx := context.Background()
	wantErr := errors.New("not found")

	ctrl := gomock.NewController(t)
	userRepo := NewMockuserRepository(ctrl)
	authRepo := NewMockauthRepository(ctrl)

	authRepo.EXPECT().
		GetAndDeleteRefreshToken(gomock.Any(), "valid-token").
		Return(domain.RefreshToken{UserID: 99, ExpiresAt: time.Now().Add(time.Hour)}, nil)

	userRepo.EXPECT().
		GetByID(gomock.Any(), int64(99)).
		Return(nil, wantErr)

	_, err := NewAuth("secret", userRepo, authRepo).LoginWithRefreshToken(ctx, "valid-token")
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
