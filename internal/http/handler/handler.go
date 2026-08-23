package handler

import (
	"github.com/nipalab/nipa/internal/usecase"
)

type usecaseContainer interface {
	AuthUsecase() *usecase.Auth
	UserUsecase() *usecase.User
}

type Handler struct {
	useCase usecaseContainer
}

func NewHandler(useCase usecaseContainer) *Handler {
	return &Handler{
		useCase: useCase,
	}
}
