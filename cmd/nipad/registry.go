package main

import "github.com/nipalab/nipa/internal/usecase"

type Registry struct {
	authUsecase   *usecase.Auth
	userUsecase   *usecase.User
	branchUsecase *usecase.Branch
}

func (r *Registry) Auth() *usecase.Auth {
	return r.authUsecase
}

func (r *Registry) User() *usecase.User {
	return r.userUsecase
}

func (r *Registry) Branch() *usecase.Branch {
	return r.branchUsecase
}
