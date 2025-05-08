package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        string
	UserID    uuid.UUID
	ExpiresAt time.Time
	CreatedAt time.Time
}
