package main

import "github.com/nipalab/nipa/internal/client/usecase"

type Registry struct {
	authUsecase *usecase.Auth
	repoUsecase *usecase.Repo
}

func (r *Registry) Auth() *usecase.Auth {
	return r.authUsecase
}

func (r *Registry) Repo() *usecase.Repo {
	return r.repoUsecase
}
