package interfaces

import (
	"music-service/internal/models"

	"github.com/google/uuid"
)

type AlbumRepository interface {
	FindByID(id uuid.UUID) (*models.Album, error)
	Save(album *models.Album) error
	Delete(id uuid.UUID) error
	GetTracks(albumID uuid.UUID) ([]*models.Track, error)
	AddTrackToAlbum(albumID, trackID uuid.UUID) error
	RemoveTrackFromAlbum(albumID, trackID uuid.UUID) error
	ListAll() ([]*models.Album, error)
}
