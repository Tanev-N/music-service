package interfaces

import (
	"music-service/internal/models"

	"github.com/google/uuid"
)

type TrackRepository interface {
	FindByID(id uuid.UUID) (*models.Track, error)
	Save(track *models.Track) error
	Delete(id uuid.UUID) error
	Search(query string) ([]*models.Track, error)
	IncrementPlayCount(trackID uuid.UUID) error
}
