package interfaces

import (
	"music-service/internal/models"

	"github.com/google/uuid"
)

type PlaylistUseCase interface {
	CreatePlaylist(userID uuid.UUID, name, description string, coverURL string) (*models.Playlist, error)
	AddTrackToPlaylist(playlistID, trackID uuid.UUID) error
	RemoveTrackFromPlaylist(playlistID, trackID uuid.UUID) error
	EditPlaylistInfo(playlistID uuid.UUID, name, description string) error
	GetPlaylistTracks(playlistID uuid.UUID) ([]*models.Track, error)
	GetPlaylistWithTracks(playlistID uuid.UUID) (*models.PlaylistTrack, error)
	GetUserPlaylists(userID uuid.UUID) ([]*models.Playlist, error)
	DeletePlaylist(playlistID, userID uuid.UUID) error
}
