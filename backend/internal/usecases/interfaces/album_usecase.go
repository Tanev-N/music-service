package interfaces

import (
	"music-service/internal/models"
	"time"

	"github.com/google/uuid"
)

type AlbumUseCase interface {
	CreateAlbum(title string, artist string, releaseDate time.Time, coverURL string) (*models.Album, error)
	AddTrackToAlbum(albumID uuid.UUID, trackID uuid.UUID) error
	RemoveTrackFromAlbum(albumID uuid.UUID, trackID uuid.UUID) error
	GetAlbumDetails(albumID uuid.UUID) (*models.Album, []*models.Track, error)
	UpdateAlbumInfo(albumID uuid.UUID, title, artist, coverURL string, releaseDate time.Time) error
	ListAll() ([]*models.Album, error)
	Delete(albumID uuid.UUID) error
}
