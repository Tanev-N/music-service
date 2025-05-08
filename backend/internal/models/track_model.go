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
	Genres    []*Genre
	PlayCount int
	Duration  int
}

// TrackUploadMetadata содержит метаданные для загрузки трека
type TrackUploadMetadata struct {
	Title      string    `json:"title"`
	AlbumID    uuid.UUID `json:"album_id,omitempty"`
	ArtistName string    `json:"artist_name"`
	Duration   int       `json:"duration,omitempty"`
	CoverURL   string    `json:"cover_url,omitempty"`
}
