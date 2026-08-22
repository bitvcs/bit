package main

import "github.com/bitvcs/bit/internal/usecase"

type Registry struct {
	authUsecase *usecase.Auth
	userUsecase *usecase.User
}

func (r *Registry) AuthUsecase() *usecase.Auth {
	return r.authUsecase
}

func (r *Registry) UserUsecase() *usecase.User {
	return r.userUsecase
}
