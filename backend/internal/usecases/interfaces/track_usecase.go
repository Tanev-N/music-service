package interfaces

import (
	"io"
	"music-service/internal/models"

	"github.com/google/uuid"
)

type TrackUseCase interface {
	SearchTracks(query string) ([]*models.Track, error)
	PlayTrack(userID uuid.UUID, trackID uuid.UUID) error
	GetTrackDetails(trackID uuid.UUID) (*models.TrackDetails, error)
	UpdateTrackMetadata(trackID uuid.UUID, metadata map[string]interface{}) error
	DeleteTrack(trackID uuid.UUID) error
	UploadTrack(fileReader io.Reader, fileSize int64, metadata models.TrackUploadMetadata) (*models.Track, error)
	GetTrackFilePath(trackID uuid.UUID) (string, error)
}
