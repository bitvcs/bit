package usecase

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nipalab/nipa/internal/client/domain"
	serverDomain "github.com/nipalab/nipa/internal/domain"
)

type loginExecutor interface {
	LoginWithUsernamePassword(ctx context.Context, url, username, password string) (*domain.LoginResult, error)
	LoginWithRefreshToken(ctx context.Context, url, refreshToken string) (*domain.LoginResult, error)
}

type secureStorage interface {
	SaveToken(data *domain.LoginResult) error
	LoadToken(serverDomain string) (*domain.LoginResult, error)
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

	err = l.secureStorage.SaveToken(loginResult)
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

	err = l.secureStorage.SaveToken(loginResult)
	if err != nil {
		return err
	}

	return nil
}

func (l *Auth) IsLoggedIn(domainUrl string) (bool, error) {
	token, err := l.secureStorage.LoadToken(domainUrl)
	if err != nil {
		return false, nil
	}
	claims := serverDomain.Claims{}
	_, _, err = jwt.NewParser().ParseUnverified(token.AccessToken, &claims)
	if err != nil {
		return false, domain.NewTokenError(err.Error())
	}
	expirationTime, err := claims.GetExpirationTime()
	if err != nil {
		return false, domain.NewTokenError(err.Error())
	}
	if expirationTime.Before(time.Now().Add(10 * time.Minute)) {
		return false, nil
	}
	return false, nil
}
