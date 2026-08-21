package usecase

import (
	"context"
	"strconv"
	"time"

	"github.com/bitvcs/bit/internal/domain"
	"github.com/bitvcs/bit/internal/usecase/dto"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	tokenExpirationMinutes = 30
)

type userRepository interface {
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByID(ctx context.Context, id int64) (*domain.User, error)
}

type authRepository interface {
	SaveRefreshToken(ctx context.Context, userID int64, refreshToken string, expiresAt int64) error
	GetAndDeleteRefreshToken(ctx context.Context, refreshToken string) (domain.RefreshToken, error)
}

type Auth struct {
	secret   string
	userRepo userRepository
	authRepo authRepository
}

func NewAuth(secret string, userRepo userRepository, authRepo authRepository) *Auth {
	return &Auth{
		secret:   secret,
		userRepo: userRepo,
		authRepo: authRepo,
	}
}

func (a *Auth) LoginWithEmailPassword(ctx context.Context, email string, password string) (*dto.LoginResult, error) {
	user, err := a.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return a.generateLoginResult(ctx, user)
}

func (a *Auth) LoginWithRefreshToken(ctx context.Context, refreshToken string) (*dto.LoginResult, error) {
	storedToken, err := a.authRepo.GetAndDeleteRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	if storedToken.ExpiresAt.Before(time.Now()) {
		return nil, domain.NewErrorUser("refresh token expired")
	}

	user, err := a.userRepo.GetByID(ctx, storedToken.UserID)
	if err != nil {
		return nil, err
	}

	return a.generateLoginResult(ctx, user)
}

func (a *Auth) generateLoginResult(ctx context.Context, user *domain.User) (*dto.LoginResult, error) {
	timeStart := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(user.ID, 10),
		IssuedAt:  jwt.NewNumericDate(timeStart),
		ExpiresAt: jwt.NewNumericDate(timeStart.Add(tokenExpirationMinutes * time.Minute)),
	}
	jwt, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(a.secret))
	if err != nil {
		return nil, domain.NewErrorInternalServer(err.Error())
	}

	refreshToken := uuid.New().String()
	tokenExpiresAt := timeStart.Add(60 * 24 * time.Hour).Unix()
	err = a.authRepo.SaveRefreshToken(ctx, user.ID, refreshToken, tokenExpiresAt)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResult{
		AccessToken:  jwt,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
	}, nil
}
