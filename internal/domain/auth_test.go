package domain

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nipalab/nipa/internal/snow"
)

func TestContextWithClaim_ClaimFromContext_RoundTrip(t *testing.T) {
	claim := Claims{
		UserID:       snow.ID(42),
		IsSuperAdmin: true,
		IsAdmin:      true,
		OrgID:        []snow.ID{snow.ID(1), snow.ID(2)},
	}

	ctx := ContextWithClaim(context.Background(), claim)
	got, ok := ClaimFromContext(ctx)
	require.True(t, ok)
	require.Equal(t, claim, got)
}

func TestClaimFromContext_EmptyContext(t *testing.T) {
	got, ok := ClaimFromContext(context.Background())
	require.False(t, ok)
	require.Equal(t, Claims{}, got)
}

func TestClaimFromContext_WrongKeyType(t *testing.T) {
	type otherKey struct{}
	ctx := context.WithValue(context.Background(), otherKey{}, "not a claim")

	got, ok := ClaimFromContext(ctx)
	require.False(t, ok)
	require.Equal(t, Claims{}, got)
}

func TestContextWithClaim_OverwritesPrevious(t *testing.T) {
	first := Claims{UserID: 1}
	second := Claims{UserID: 2}

	ctx := ContextWithClaim(context.Background(), first)
	ctx = ContextWithClaim(ctx, second)

	got, ok := ClaimFromContext(ctx)
	require.True(t, ok)
	require.Equal(t, second, got)
}

func TestClaims_Fields(t *testing.T) {
	claim := Claims{
		UserID:       snow.ID(10),
		IsSuperAdmin: false,
		IsAdmin:      true,
		OrgID:        []snow.ID{snow.ID(5), snow.ID(6), snow.ID(7)},
	}

	require.Equal(t, snow.ID(10), claim.UserID)
	require.False(t, claim.IsSuperAdmin)
	require.True(t, claim.IsAdmin)
	require.Len(t, claim.OrgID, 3)
}
