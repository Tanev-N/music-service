package interfaces

import (
	"music-service/internal/models"

	"github.com/google/uuid"
)

type TrackUseCase interface {
	SearchTracks(query string) ([]*models.Track, error)
	PlayTrack(userID uuid.UUID, trackID uuid.UUID) error
	GetTrackDetails(trackID uuid.UUID) (*models.Track, error)
	UpdateTrackMetadata(trackID uuid.UUID, metadata map[string]interface{}) error
	DeleteTrack(trackID uuid.UUID) error
}
