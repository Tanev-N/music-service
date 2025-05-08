package models

import (
	"time"

	"github.com/google/uuid"
)

type Album struct {
	ID          uuid.UUID
	Title       string
	Artist      string
	ReleaseDate time.Time
	CoverURL    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
