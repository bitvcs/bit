package main

import (
	"testing"

	"github.com/bitvcs/bit/internal/usecase"
	"github.com/stretchr/testify/require"
)

func TestRegistry(t *testing.T) {
	reg := &Registry{
		authUsecase: usecase.NewAuth("secret", nil, nil),
		userUsecase: usecase.NewUser(nil),
	}

	require.Same(t, reg.authUsecase, reg.AuthUsecase())
	require.Same(t, reg.userUsecase, reg.UserUsecase())
}
