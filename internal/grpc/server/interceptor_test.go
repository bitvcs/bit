package server

import (
	"context"
	"errors"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/nipalab/nipa/internal/domain"
)

type mockTokenValidator struct {
	claims *domain.Claims
	err    error
}

func (m *mockTokenValidator) ValidateToken(_ context.Context, _ string) (*domain.Claims, error) {
	return m.claims, m.err
}

func unaryHandler(t *testing.T) grpc.UnaryHandler {
	t.Helper()
	return func(_ context.Context, _ interface{}) (interface{}, error) {
		return "ok", nil
	}
}

func TestInterceptor_PublicRoutes(t *testing.T) {
	handler := unaryHandler(t)
	publicRoutes := []string{
		"/greet.NipaService/LoginWithUsernamePassword",
		"/greet.NipaService/LoginWithRefreshToken",
		"/nipa.AuthService/health",
	}

	for _, route := range publicRoutes {
		t.Run(route, func(t *testing.T) {
			interceptor := NewInterceptor(nil)
			resp, err := interceptor.JWTUnary()(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: route}, handler)
			require.NoError(t, err)
			require.Equal(t, "ok", resp)
		})
	}
}

func TestInterceptor_MissingMetadata(t *testing.T) {
	interceptor := NewInterceptor(nil)
	handler := unaryHandler(t)

	_, err := interceptor.JWTUnary()(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: "/nipa.SomeService/Method"}, handler)
	require.Error(t, err)
	require.Equal(t, codes.Unauthenticated, status.Code(err))
	require.Contains(t, status.Convert(err).Message(), "missing metadata")
}

func TestInterceptor_MissingAuthorizationHeader(t *testing.T) {
	interceptor := NewInterceptor(nil)
	handler := unaryHandler(t)

	ctx := metadata.NewIncomingContext(context.Background(), metadata.MD{})
	_, err := interceptor.JWTUnary()(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/nipa.SomeService/Method"}, handler)
	require.Error(t, err)
	require.Equal(t, codes.Unauthenticated, status.Code(err))
	require.Contains(t, status.Convert(err).Message(), "missing authorization header")
}

func TestInterceptor_InvalidAuthorizationHeader(t *testing.T) {
	interceptor := NewInterceptor(nil)
	handler := unaryHandler(t)

	tests := []struct {
		name  string
		value string
	}{
		{"no bearer prefix", "token123"},
		{"lowercase bearer", "bearer token123"},
		{"only bearer with no token", "Bearer"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := metadata.NewIncomingContext(context.Background(), metadata.MD{
				"authorization": []string{tt.value},
			})
			_, err := interceptor.JWTUnary()(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/nipa.SomeService/Method"}, handler)
			require.Error(t, err)
			require.Equal(t, codes.Unauthenticated, status.Code(err))
			require.Contains(t, status.Convert(err).Message(), "invalid authorization header")
		})
	}
}

func TestInterceptor_InvalidToken(t *testing.T) {
	validator := &mockTokenValidator{err: errors.New("bad token")}
	interceptor := NewInterceptor(validator)
	handler := unaryHandler(t)

	ctx := metadata.NewIncomingContext(context.Background(), metadata.MD{
		"authorization": []string{"Bearer bad-token"},
	})
	_, err := interceptor.JWTUnary()(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/nipa.SomeService/Method"}, handler)
	require.Error(t, err)
	require.Equal(t, codes.Unauthenticated, status.Code(err))
	require.Contains(t, status.Convert(err).Message(), "invalid token")
}

func TestInterceptor_ValidToken(t *testing.T) {
	wantClaims := &domain.Claims{
		RegisteredClaims: jwt.RegisteredClaims{Subject: "42"},
		UserID:           42,
		IsAdmin:          true,
	}
	validator := &mockTokenValidator{claims: wantClaims}
	interceptor := NewInterceptor(validator)

	var gotClaims domain.Claims
	handler := func(ctx context.Context, _ interface{}) (interface{}, error) {
		claims, ok := domain.ClaimFromContext(ctx)
		require.True(t, ok)
		gotClaims = claims
		return "ok", nil
	}

	ctx := metadata.NewIncomingContext(context.Background(), metadata.MD{
		"authorization": []string{"Bearer valid-token"},
	})
	resp, err := interceptor.JWTUnary()(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/nipa.SomeService/Method"}, handler)
	require.NoError(t, err)
	require.Equal(t, "ok", resp)
	require.Equal(t, wantClaims.UserID, gotClaims.UserID)
	require.True(t, gotClaims.IsAdmin)
}
