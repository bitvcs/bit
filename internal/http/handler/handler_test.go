package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	database "github.com/bitvcs/bit/db"
	"github.com/bitvcs/bit/internal/domain"
	httpApp "github.com/bitvcs/bit/internal/http"
	"github.com/bitvcs/bit/internal/http/model"
	sqlcSqlite "github.com/bitvcs/bit/internal/repository/sqlc/sqlite"
	"github.com/bitvcs/bit/internal/repository/sqlite"
	"github.com/bitvcs/bit/internal/snow"
	"github.com/bitvcs/bit/internal/usecase"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

type fakeAppContext struct {
	body       []byte
	statusCode int
	response   any
}

func (f *fakeAppContext) Context() context.Context { return context.Background() }
func (f *fakeAppContext) Claims() *domain.Claims   { return nil }

func (f *fakeAppContext) ReadJson(v any) error {
	return json.Unmarshal(f.body, v)
}

func (f *fakeAppContext) WriteJson(statusCode int, v any) error {
	f.statusCode = statusCode
	f.response = v
	return nil
}

func (f *fakeAppContext) HandleError(err error) {
	apiErr, ok := err.(*domain.Error)
	if ok {
		f.statusCode = apiErr.Code
		f.response = model.NewAPIError(apiErr.Message)
		return
	}
	f.statusCode = http.StatusInternalServerError
	f.response = model.NewAPIError("internal unknown error")
}

type handlerRegistry struct {
	auth *usecase.Auth
	user *usecase.User
}

func (r *handlerRegistry) AuthUsecase() *usecase.Auth { return r.auth }
func (r *handlerRegistry) UserUsecase() *usecase.User { return r.user }

func newHandlerTestSetup(t *testing.T) (*Handler, *sqlite.Auth, snow.ID) {
	t.Helper()

	dbConn, err := database.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, dbConn.Close()) })
	dbConn.SetMaxOpenConns(1)
	require.NoError(t, database.MigrateUp(dbConn, "sqlite3"))

	q := sqlcSqlite.New(dbConn)
	userID, err := q.UserCreate(context.Background(), sqlcSqlite.UserCreateParams{
		Name:     "alice",
		Email:    "alice@example.com",
		Password: "hashed-password",
		IsAdmin:  true,
	})
	require.NoError(t, err)

	authRepo := sqlite.NewAuthRepository(dbConn)
	node, _ := snow.NewNode(1)
	reg := &handlerRegistry{
		auth: usecase.NewAuth("test-secret", sqlite.NewUserRepository(dbConn), authRepo),
		user: usecase.NewUser(node),
	}
	return NewHandler(reg), authRepo, snow.ID(userID)
}

func TestNewHandler(t *testing.T) {
	h, _, _ := newHandlerTestSetup(t)
	require.NotNil(t, h)
}

func TestHandler_AuthRefreshToken(t *testing.T) {
	h, authRepo, userID := newHandlerTestSetup(t)

	require.NoError(t, authRepo.SaveRefreshToken(
		context.Background(), userID, "valid-refresh-token", time.Now().Add(time.Hour)))

	appCtx := &fakeAppContext{body: []byte(`{"refresh_token":"valid-refresh-token"}`)}
	h.AuthRefreshToken(appCtx)

	require.Equal(t, http.StatusOK, appCtx.statusCode)

	loginResponse, ok := appCtx.response.(model.LoginResponse)
	require.True(t, ok)
	require.Equal(t, "Bearer", loginResponse.TokenType)
	require.NotEmpty(t, loginResponse.AccessToken)
	require.NotEmpty(t, loginResponse.RefreshToken)
}

func TestHandler_AuthRefreshToken_Expired(t *testing.T) {
	h, authRepo, userID := newHandlerTestSetup(t)

	require.NoError(t, authRepo.SaveRefreshToken(
		context.Background(), userID, "expired-refresh-token", time.Now().Add(-time.Hour)))

	appCtx := &fakeAppContext{body: []byte(`{"refresh_token":"expired-refresh-token"}`)}
	h.AuthRefreshToken(appCtx)

	require.Equal(t, http.StatusBadRequest, appCtx.statusCode)

	apiErr, ok := appCtx.response.(*model.APIError)
	require.True(t, ok)
	require.Equal(t, "refresh token expired", apiErr.Error)
}

func TestHandler_AuthRefreshToken_UnknownToken(t *testing.T) {
	h, _, _ := newHandlerTestSetup(t)

	appCtx := &fakeAppContext{body: []byte(`{"refresh_token":"unknown"}`)}
	h.AuthRefreshToken(appCtx)

	require.Equal(t, http.StatusNotFound, appCtx.statusCode)

	apiErr, ok := appCtx.response.(*model.APIError)
	require.True(t, ok)
	require.Equal(t, "record not found", apiErr.Error)
}

func TestHandler_AuthRefreshToken_InvalidJson(t *testing.T) {
	h, _, _ := newHandlerTestSetup(t)

	appCtx := &fakeAppContext{body: []byte(`{"refresh_token":`)}
	h.AuthRefreshToken(appCtx)

	require.Equal(t, http.StatusInternalServerError, appCtx.statusCode)
}

func TestHandler_AuthLogin(t *testing.T) {
	h, _, userID := newHandlerTestSetup(t)

	appCtx := &fakeAppContext{body: []byte(`{"email":"alice@example.com","password":"whatever"}`)}
	h.AuthLogin(appCtx)

	require.Equal(t, http.StatusOK, appCtx.statusCode)

	loginResponse, ok := appCtx.response.(model.LoginResponse)
	require.True(t, ok)
	require.Equal(t, "Bearer", loginResponse.TokenType)
	require.NotEmpty(t, loginResponse.AccessToken)
	require.NotEmpty(t, loginResponse.RefreshToken)

	claims := parseHandlerAccessToken(t, loginResponse.AccessToken, "test-secret")
	require.Equal(t, userID, claims.UserID)
}

func TestHandler_AuthLogin_UnknownUser(t *testing.T) {
	h, _, _ := newHandlerTestSetup(t)

	appCtx := &fakeAppContext{body: []byte(`{"email":"ghost@example.com","password":"whatever"}`)}
	h.AuthLogin(appCtx)

	require.Equal(t, http.StatusBadRequest, appCtx.statusCode)

	apiErr, ok := appCtx.response.(*model.APIError)
	require.True(t, ok)
	require.Equal(t, "invalid email or password", apiErr.Error)
}

func TestHandler_AuthLogin_InvalidJson(t *testing.T) {
	h, _, _ := newHandlerTestSetup(t)

	appCtx := &fakeAppContext{body: []byte(`{"email":`)}
	h.AuthLogin(appCtx)

	require.Equal(t, http.StatusInternalServerError, appCtx.statusCode)

	apiErr, ok := appCtx.response.(*model.APIError)
	require.True(t, ok)
	require.Equal(t, "internal unknown error", apiErr.Error)
}

var _ httpApp.AppContext = (*fakeAppContext)(nil)

func parseHandlerAccessToken(t *testing.T, token string, secret string) *domain.Claims {
	t.Helper()

	claims := domain.Claims{}
	parsed, err := jwt.ParseWithClaims(token, &claims, func(_ *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	require.NoError(t, err)
	require.True(t, parsed.Valid)
	return &claims
}
