package domain

import (
	"time"

	"github.com/nipalab/nipa/internal/snow"
)

type Organization struct {
	ID        snow.ID    `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Deleted   bool       `json:"deleted"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	Name      string     `json:"name"`
	Slug      string     `json:"slug"`
}

type Project struct {
	ID          snow.ID    `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Deleted     bool       `json:"deleted"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Slug        string     `json:"slug"`
	OrgID       snow.ID    `json:"org_id"`
}
