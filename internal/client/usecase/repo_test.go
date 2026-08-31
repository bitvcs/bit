package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"

	"github.com/nipalab/nipa/internal/client/domain"
	serverDomain "github.com/nipalab/nipa/internal/domain"
	"github.com/nipalab/nipa/internal/snow"
)

func TestNewRepo(t *testing.T) {
	auth := NewAuth(nil, nil, nil)
	repo := NewRepo(auth)
	require.Equal(t, auth, repo.auth)
}

func TestRepo_Clone_Success(t *testing.T) {
	token := signTestToken(t, "secret")
	storage := &stubSecureStorage{loadResult: &domain.LoginResult{AccessToken: token}}
	auth := NewAuth(nil, storage, nil)
	repo := NewRepo(auth)

	err := repo.Clone(context.Background(), "example.com", "org", "project", "main", "/src", "/target")
	require.NoError(t, err)
}

func TestRepo_Clone_Success_NeedLogin(t *testing.T) {
	token := signTestToken(t, "secret")
	storage := &stubSecureStorage{loadErr: errors.New("not found")}
	input := &stubUserInput{username: "apin", password: "secret"}
	executor := &stubLoginExecutor{usernameResult: &domain.LoginResult{
		AccessToken: token,
		Host:        "example.com",
	}}
	auth := NewAuth(executor, storage, input)
	repo := NewRepo(auth)

	err := repo.Clone(context.Background(), "example.com", "org", "project", "main", "/src", "/target")
	require.NoError(t, err)
	require.Equal(t, "apin", executor.lastUsername)
	require.Equal(t, "secret", executor.lastPassword)
}

func TestRepo_Clone_Error_PromptFailed(t *testing.T) {
	wantErr := errors.New("prompt interrupted")
	storage := &stubSecureStorage{loadErr: errors.New("not found")}
	input := &stubUserInput{err: wantErr}
	auth := NewAuth(nil, storage, input)
	repo := NewRepo(auth)

	err := repo.Clone(context.Background(), "example.com", "org", "project", "main", "/src", "/target")
	require.ErrorIs(t, err, wantErr)
}

func TestRepo_Clone_Error_LoginFailed(t *testing.T) {
	wantErr := errors.New("login failed")
	storage := &stubSecureStorage{loadErr: errors.New("not found")}
	input := &stubUserInput{username: "apin", password: "secret"}
	executor := &stubLoginExecutor{usernameErr: wantErr}
	auth := NewAuth(executor, storage, input)
	repo := NewRepo(auth)

	err := repo.Clone(context.Background(), "example.com", "org", "project", "main", "/src", "/target")
	require.ErrorIs(t, err, wantErr)
}

func TestRepo_Clone_Error_ExpiredToken(t *testing.T) {
	claims := serverDomain.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "42",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
		},
		UserID: snow.ID(42),
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("secret"))
	require.NoError(t, err)

	storage := &stubSecureStorage{loadResult: &domain.LoginResult{AccessToken: token}}
	input := &stubUserInput{username: "apin", password: "secret"}
	executor := &stubLoginExecutor{usernameResult: &domain.LoginResult{
		AccessToken: "fresh-token",
		Host:        "example.com",
	}}
	auth := NewAuth(executor, storage, input)
	repo := NewRepo(auth)

	err = repo.Clone(context.Background(), "example.com", "org", "project", "main", "/src", "/target")
	require.NoError(t, err)
	require.Equal(t, "apin", executor.lastUsername)
	require.Equal(t, "secret", executor.lastPassword)
}

func TestRepo_Clone_MalformedToken(t *testing.T) {
	storage := &stubSecureStorage{loadResult: &domain.LoginResult{AccessToken: "not-a-jwt"}}
	input := &stubUserInput{username: "apin", password: "secret"}
	executor := &stubLoginExecutor{usernameResult: &domain.LoginResult{
		AccessToken: "fresh-token",
		Host:        "example.com",
	}}
	auth := NewAuth(executor, storage, input)
	repo := NewRepo(auth)

	err := repo.Clone(context.Background(), "example.com", "org", "project", "main", "/src", "/target")
	require.NoError(t, err)
	require.Equal(t, "apin", executor.lastUsername)
	require.Equal(t, "secret", executor.lastPassword)
}
