package securestorage

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"

	"github.com/nipalab/nipa/internal/client/domain"
)

func TestStorage_SaveAndLoadToken(t *testing.T) {
	keyring.MockInit()
	s := New()

	data := &domain.LoginResult{
		Host:         "example.com",
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		ExpiresIn:    1800,
	}

	require.NoError(t, s.SaveToken(data))

	got, err := s.LoadToken("example.com")
	require.NoError(t, err)
	require.Equal(t, data, got)
}

func TestStorage_SaveToken_KeyringError(t *testing.T) {
	keyring.MockInitWithError(errors.New("keyring failure"))
	s := New()

	err := s.SaveToken(&domain.LoginResult{
		Host:         "example.com",
		AccessToken:  "access",
		RefreshToken: "refresh",
	})
	require.Error(t, err)
}

func TestStorage_LoadToken_NotFound(t *testing.T) {
	keyring.MockInit()
	s := New()

	_, err := s.LoadToken("unknown.example.com")
	require.Error(t, err)
	require.Equal(t, keyring.ErrNotFound, err)
}

func TestStorage_LoadToken_KeyringError(t *testing.T) {
	keyring.MockInitWithError(errors.New("keyring failure"))
	s := New()

	_, err := s.LoadToken("example.com")
	require.Error(t, err)
}

func TestStorage_LoadToken_InvalidJSON(t *testing.T) {
	keyring.MockInit()
	s := New()

	require.NoError(t, keyring.Set(serviceName, "example.com", "not-valid-json"))
	_, err := s.LoadToken("example.com")
	require.Error(t, err)
}

func TestStorage_LoadToken_Empty(t *testing.T) {
	keyring.MockInit()
	s := New()

	require.NoError(t, keyring.Set(serviceName, "example.com", "{}"))

	got, err := s.LoadToken("example.com")
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Empty(t, got.AccessToken)
}
