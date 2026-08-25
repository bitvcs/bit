package usecase

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nipalab/nipa/internal/domain"
	"github.com/nipalab/nipa/internal/snow"
	"github.com/nipalab/nipa/internal/usecase/dto"
)

const (
	tokenExpirationMinutes = 30
)

//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=auth_mock_test.go -package=usecase
type userRepository interface {
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByID(ctx context.Context, id snow.ID) (*domain.User, error)
}

//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=auth_mock_test.go -package=usecase
type authRepository interface {
	SaveRefreshToken(ctx context.Context, userID snow.ID, refreshToken string, expiresAt time.Time) error
	GetAndDeleteRefreshToken(ctx context.Context, refreshToken string) (*domain.RefreshToken, error)
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
		if domain.IsErrorNotFound(err) {
			return nil, domain.NewErrorUser("invalid email or password")
		}
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

func (a *Auth) ValidateToken(ctx context.Context, tokenString string) (*domain.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.NewErrorUser("unexpected signing method")
		}
		return []byte(a.secret), nil
	})

	if err != nil {
		return nil, domain.NewErrorUser("invalid token")
	}

	if claims, ok := token.Claims.(*domain.Claims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, domain.NewErrorUser("invalid token")
	}
}

func (a *Auth) generateLoginResult(ctx context.Context, user *domain.User) (*dto.LoginResult, error) {
	timeStart := time.Now()
	claims := domain.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID.Base36(),
			IssuedAt:  jwt.NewNumericDate(timeStart),
			ExpiresAt: jwt.NewNumericDate(timeStart.Add(tokenExpirationMinutes * time.Minute)),
		},
		UserID:       user.ID,
		IsSuperAdmin: user.IsSuperAdmin,
		IsAdmin:      user.IsAdmin,
		OrgID:        []snow.ID{},
	}
	jwt, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(a.secret))
	if err != nil {
		return nil, domain.NewErrorInternalServer(err.Error())
	}

	refreshToken := uuid.New().String()
	tokenExpiresAt := timeStart.Add(60 * 24 * time.Hour)
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

func (a *Auth) HasProjectAccess(ctx context.Context, projectID snow.ID, permission domain.Permission) bool {
	claim, ok := domain.ClaimFromContext(ctx)
	if !ok {
		return false
	}

	if claim.IsSuperAdmin || claim.IsAdmin {
		return true
	}

	//TODO: check PBAC rules for the user and projectID with the required permission

	return false
}
