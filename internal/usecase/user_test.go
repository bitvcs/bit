package usecase

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewUser(t *testing.T) {
	require.NotNil(t, NewUser())
}
