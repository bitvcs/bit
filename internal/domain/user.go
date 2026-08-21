package domain

import "time"

type User struct {
	ID           int64      `json:"id"`
	Name         string     `json:"name"`
	Email        string     `json:"email"`
	Password     string     `json:"password"`
	PhotoUrl     string     `json:"photo_url"`
	IsSuperAdmin bool       `json:"is_super_admin"`
	IsAdmin      bool       `json:"is_admin"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	Deleted      bool       `json:"deleted"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

type Group struct {
	ID          int64      `json:"id"`
	OrgID       int64      `json:"org_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Deleted     bool       `json:"deleted"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	MemberIDs   []int64    `json:"member_ids"`
}

type RefreshToken struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}
