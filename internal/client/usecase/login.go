package usecase

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nipalab/nipa/internal/client/domain"
	serverDomain "github.com/nipalab/nipa/internal/domain"
)

type userInput interface {
	PromptUsernameAndPassword() (username, password string, err error)
}

type loginExecutor interface {
	LoginWithUsernamePassword(ctx context.Context, host, username, password string) (*domain.LoginResult, error)
	LoginWithRefreshToken(ctx context.Context, host, refreshToken string) (*domain.LoginResult, error)
}

type secureStorage interface {
	SaveToken(data *domain.LoginResult) error
	LoadToken(host string) (*domain.LoginResult, error)
}

type Auth struct {
	loginExecutor loginExecutor
	secureStorage secureStorage
	userInput     userInput
}

func NewAuth(executor loginExecutor, storage secureStorage, input userInput) *Auth {
	return &Auth{
		loginExecutor: executor,
		secureStorage: storage,
		userInput:     input,
	}
}

func (a *Auth) SetLoginExecutor(executor loginExecutor) {
	a.loginExecutor = executor
}

func (l *Auth) LoginWithUsernamePassword(ctx context.Context, host, username, password string) error {
	if username == "" {
		return domain.NewUserError("username cannot be empty")
	}

	loginResult, err := l.loginExecutor.LoginWithUsernamePassword(ctx, host, username, password)
	if err != nil {
		return err
	}

	err = l.secureStorage.SaveToken(loginResult)
	if err != nil {
		return err
	}

	return nil
}

func (l *Auth) LoginWithRefreshToken(ctx context.Context, host, refreshToken string) error {
	if refreshToken == "" {
		return domain.NewUserError("refresh token cannot be empty")
	}

	loginResult, err := l.loginExecutor.LoginWithRefreshToken(ctx, host, refreshToken)
	if err != nil {
		return err
	}

	err = l.secureStorage.SaveToken(loginResult)
	if err != nil {
		return err
	}

	return nil
}

func (l *Auth) isLoggedIn(host string) (bool, error) {
	token, err := l.secureStorage.LoadToken(host)
	if err != nil {
		return false, nil
	}
	claims := serverDomain.Claims{}
	_, _, err = jwt.NewParser().ParseUnverified(token.AccessToken, &claims)
	if err != nil {
		return false, domain.NewTokenError(err.Error())
	}
	return true, nil
}

func (l *Auth) MakeSureLoggedIn(ctx context.Context, host string) error {
	if loggedIn, err := l.isLoggedIn(host); !loggedIn || err != nil {
		username, password, err := l.userInput.PromptUsernameAndPassword()
		if err != nil {
			return err
		}
		err = l.LoginWithUsernamePassword(ctx, host, username, password)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *Auth) GetToken(ctx context.Context, host string) (string, error) {
	loginResult, err := l.secureStorage.LoadToken(host)
	if err != nil {
		return "", err
	}
	return loginResult.AccessToken, nil
}
