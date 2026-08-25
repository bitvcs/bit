package main

import "github.com/nipalab/nipa/internal/usecase"

type Registry struct {
	authUsecase *usecase.Auth
	userUsecase *usecase.User
}

func (r *Registry) Auth() *usecase.Auth {
	return r.authUsecase
}

func (r *Registry) User() *usecase.User {
	return r.userUsecase
}
