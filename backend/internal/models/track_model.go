package models

import (
	"time"

	"github.com/google/uuid"
)

type Track struct {
	ID         uuid.UUID
	Title      string
	Duration   int
	FilePath   string
	AlbumID    uuid.UUID
	ArtistName string
	CoverURL   string
	AddedDate  time.Time
	UpdatedAt  time.Time
	PlayCount  int
}

type TrackDetails struct {
	Track
	Album     *Album
	PlayCount int
	Duration  int
}
