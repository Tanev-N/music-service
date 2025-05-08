package interfaces

import (
	"music-service/internal/models"

	"github.com/google/uuid"
)

type PlaylistRepository interface {
	FindByID(id uuid.UUID) (*models.Playlist, error)
	Save(playlist *models.Playlist) error
	Delete(id uuid.UUID) error
	AddTrack(playlistID uuid.UUID, trackID uuid.UUID) error
	RemoveTrack(playlistID uuid.UUID, trackID uuid.UUID) error
	GetTracks(playlistID uuid.UUID) ([]*models.Track, error)
	GetUserPlaylists(userID uuid.UUID) ([]*models.Playlist, error)
}
