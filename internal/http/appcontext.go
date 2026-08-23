package http

import (
	"context"

	"github.com/nipalab/nipa/internal/domain"
)

type AppContext interface {
	Context() context.Context
	Claims() *domain.Claims
	ReadJson(v any) error
	WriteJson(statusCode int, v any) error
	HandleError(err error)
}
