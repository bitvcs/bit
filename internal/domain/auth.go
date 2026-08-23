package domain

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/nipalab/nipa/internal/snow"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID       snow.ID   `json:"user_id"`
	IsSuperAdmin bool      `json:"is_superadmin"`
	IsAdmin      bool      `json:"is_admin"`
	OrgID        []snow.ID `json:"org_id"`
}
