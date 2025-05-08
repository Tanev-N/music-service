package interfaces

import (
	"io"
	"music-service/internal/models"

	"github.com/google/uuid"
)

type TrackRepository interface {
	FindByID(id uuid.UUID) (*models.Track, error)
	Save(track *models.Track) error
	Delete(id uuid.UUID) error
	Search(query string) ([]*models.Track, error)
	IncrementPlayCount(trackID uuid.UUID) error
	SaveTrackFile(trackID uuid.UUID, fileReader io.Reader, fileSize int64) (string, error)
	GetStorageDir() string
	GetGenresForTrack(trackID uuid.UUID) ([]*models.Genre, error)
}
