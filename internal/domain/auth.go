package domain

import (
	"github.com/bitvcs/bit/internal/snow"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID       snow.ID   `json:"user_id"`
	IsSuperAdmin bool      `json:"is_superadmin"`
	IsAdmin      bool      `json:"is_admin"`
	OrgID        []snow.ID `json:"org_id"`
}
