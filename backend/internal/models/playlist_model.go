package models

import (
	"time"

	"github.com/google/uuid"
)

type Playlist struct {
	ID          uuid.UUID
	Name        string
	Description string
	UserID      uuid.UUID
	CoverURL    string
	CreatedDate time.Time
	UpdatedAt   time.Time
}

type PlaylistTrack struct {
	Playlist Playlist
	Tracks   []*Track
}
