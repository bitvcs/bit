package domain

import (
	"time"

	"github.com/nipalab/nipa/internal/snow"
)

type Branch struct {
	ID        snow.ID   `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}
