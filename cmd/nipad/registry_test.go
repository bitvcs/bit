package main

import (
	"testing"

	"github.com/nipalab/nipa/internal/usecase"
	"github.com/stretchr/testify/require"
)

func TestRegistry(t *testing.T) {
	reg := &Registry{
		authUsecase: usecase.NewAuth("secret", nil, nil),
		userUsecase: usecase.NewUser(nil),
	}

	require.Same(t, reg.authUsecase, reg.Auth())
	require.Same(t, reg.userUsecase, reg.User())
}
