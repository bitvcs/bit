package main

import "github.com/nipalab/nipa/internal/client/usecase"

type Registry struct {
	authUsecase *usecase.Auth
}

func (r *Registry) Auth() *usecase.Auth {
	return r.authUsecase
}
