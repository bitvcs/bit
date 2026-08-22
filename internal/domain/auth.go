package domain

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	jwt.RegisteredClaims
	UserID       int64   `json:"user_id"`
	IsSuperAdmin bool    `json:"is_superadmin"`
	IsAdmin      bool    `json:"is_admin"`
	OrgID        []int64 `json:"org_id"`
}
