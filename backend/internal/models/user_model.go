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
	ID         uuid.UUID
	Login      string
	Password   string
	Permission Permission
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (p Permission) IsValid() bool {
	switch p {
	case UserPermission, AdminPermission:
		return true
	default:
		return false
	}
}
