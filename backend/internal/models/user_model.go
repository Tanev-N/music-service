package models

import (
	"time"

	"github.com/google/uuid"
)

type Permission string

const (
	UserPermission      Permission = "user"
	AdminPermission     Permission = "admin"
	ModeratorPermission Permission = "moderator"
)

type User struct {
	ID         uuid.UUID  `json:"id"`
	Login      string     `json:"login"`
	Password   string     `json:"-" yaml:"-" xml:"-" bson:"-"`
	Permission Permission `json:"permission"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

func (p Permission) IsValid() bool {
	switch p {
	case UserPermission, AdminPermission, ModeratorPermission:
		return true
	default:
		return false
	}
}
