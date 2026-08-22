package model

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewAPIError(t *testing.T) {
	err := NewAPIError("something failed")
	require.Equal(t, &APIError{Error: "something failed"}, err)
}
