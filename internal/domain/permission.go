package domain

import (
	"time"

	"github.com/nipalab/nipa/internal/snow"
)

const (
	PermissionRead  = 1 << iota // 1
	PermissionWrite             // 2
	PermissionLock              // 4
	PermissionAdmin = 1 << 16   // 65536
)

type PBACRule struct {
	ID          int64     `json:"id"`
	UserID      *int64    `json:"user_id,omitempty"`
	GroupID     *int64    `json:"group_id,omitempty"`
	OrgID       int64     `json:"org_id"`
	ProjectID   int64     `json:"project_id"`
	PathPattern string    `json:"path_pattern"`
	Permission  int       `json:"permission"`
	CreatedAt   time.Time `json:"created_at"`
}

type ProjectPathPermission struct {
	ID          int64     `json:"id"`
	ProjectID   snow.ID   `json:"project_id"`
	PathPattern string    `json:"path_pattern"`
	IsAllowed   bool      `json:"is_allowed"`
	CreatedAt   time.Time `json:"created_at"`
}
