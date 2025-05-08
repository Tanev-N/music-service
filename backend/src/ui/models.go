package ui

import (
	"music-service/internal/models"

	"github.com/google/uuid"
)

type PlaylistTrack struct {
	ID          uuid.UUID
	Name        string
	UserID      uuid.UUID
	Description string
	Tracks      []*models.Track
}
