package usecase

import (
	"context"

	"github.com/nipalab/nipa/internal/client/domain"
)

type loginExecutor interface {
	LoginWithUsernamePassword(ctx context.Context, url, username, password string) (*domain.LoginResult, error)
	LoginWithRefreshToken(ctx context.Context, url, refreshToken string) (*domain.LoginResult, error)
}

type secureStorage interface {
	SaveToken(domain, accessToken, refreshToken string, expiresIn int) error
}

type Auth struct {
	loginExecutor loginExecutor
	secureStorage secureStorage
}

func NewAuth(executor loginExecutor, storage secureStorage) *Auth {
	return &Auth{
		loginExecutor: executor,
		secureStorage: storage,
	}
}

func (l *Auth) LoginWithUsernamePassword(ctx context.Context, url, username, password string) error {
	if username == "" {
		return domain.NewUserError("username cannot be empty")
	}

	loginResult, err := l.loginExecutor.LoginWithUsernamePassword(ctx, url, username, password)
	if err != nil {
		return err
	}

	err = l.secureStorage.SaveToken(loginResult.Domain, loginResult.AccessToken, loginResult.RefreshToken, loginResult.ExpiresIn)
	if err != nil {
		return err
	}

	return nil
}

func (l *Auth) LoginWithRefreshToken(ctx context.Context, url, refreshToken string) error {
	if refreshToken == "" {
		return domain.NewUserError("refresh token cannot be empty")
	}

	loginResult, err := l.loginExecutor.LoginWithRefreshToken(ctx, url, refreshToken)
	if err != nil {
		return err
	}

	err = l.secureStorage.SaveToken(loginResult.Domain, loginResult.AccessToken, loginResult.RefreshToken, loginResult.ExpiresIn)
	if err != nil {
		return err
	}

	return nil
}
