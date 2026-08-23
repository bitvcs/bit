package usecase

import "github.com/nipalab/nipa/internal/snow"

type User struct {
	snowNode snow.Node
}

func NewUser(snowNode snow.Node) *User {
	return &User{snowNode: snowNode}
}
