package usecases

import (
	"errors"
	"fmt"
	"music-service/internal/models"
	"music-service/internal/repository/interfaces"
	usecaseInterfaces "music-service/internal/usecases/interfaces"
	"time"

	"github.com/google/uuid"
)

type playlistUseCase struct {
	playlistRepo interfaces.PlaylistRepository
	trackRepo    interfaces.TrackRepository
	userRepo     interfaces.UserRepository
}

func NewPlaylistUseCase(
	playlistRepo interfaces.PlaylistRepository,
	trackRepo interfaces.TrackRepository,
	userRepo interfaces.UserRepository,
) usecaseInterfaces.PlaylistUseCase {
	return &playlistUseCase{
		playlistRepo: playlistRepo,
		trackRepo:    trackRepo,
		userRepo:     userRepo,
	}
}

func (uc *playlistUseCase) CreatePlaylist(userID uuid.UUID, name, description string, coverURL string) (*models.Playlist, error) {
	if _, err := uc.userRepo.FindByID(userID); err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if len(name) < 2 {
		return nil, errors.New("playlist name must be at least 2 characters")
	}
	if len(name) > 100 {
		return nil, errors.New("playlist name is too long")
	}
	if len(description) > 500 {
		return nil, errors.New("description is too long")
	}

	playlist := &models.Playlist{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		UserID:      userID,
		CoverURL:    coverURL,
		CreatedDate: time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := uc.playlistRepo.Save(playlist); err != nil {
		return nil, fmt.Errorf("failed to create playlist: %w", err)
	}

	return playlist, nil
}

func (uc *playlistUseCase) AddTrackToPlaylist(playlistID, trackID uuid.UUID) error {
	playlist, err := uc.playlistRepo.FindByID(playlistID)
	if err != nil {
		return fmt.Errorf("playlist not found: %w", err)
	}

	if _, err := uc.trackRepo.FindByID(trackID); err != nil {
		return fmt.Errorf("track not found: %w", err)
	}

	tracks, err := uc.playlistRepo.GetTracks(playlistID)
	if err != nil {
		return fmt.Errorf("failed to get playlist tracks: %w", err)
	}

	for _, t := range tracks {
		if t.ID == trackID {
			return errors.New("track already exists in playlist")
		}
	}

	if len(tracks) >= 1000 {
		return errors.New("playlist track limit reached")
	}

	if err := uc.playlistRepo.AddTrack(playlistID, trackID); err != nil {
		return fmt.Errorf("failed to add track: %w", err)
	}

	playlist.UpdatedAt = time.Now()
	return uc.playlistRepo.Save(playlist)
}

func (uc *playlistUseCase) RemoveTrackFromPlaylist(playlistID, trackID uuid.UUID) error {
	playlist, err := uc.playlistRepo.FindByID(playlistID)
	if err != nil {
		return fmt.Errorf("playlist not found: %w", err)
	}

	tracks, err := uc.playlistRepo.GetTracks(playlistID)
	if err != nil {
		return fmt.Errorf("failed to get playlist tracks: %w", err)
	}

	found := false
	for _, t := range tracks {
		if t.ID == trackID {
			found = true
			break
		}
	}

	if !found {
		return errors.New("track not found in playlist")
	}

	if err := uc.playlistRepo.RemoveTrack(playlistID, trackID); err != nil {
		return fmt.Errorf("failed to remove track: %w", err)
	}

	playlist.UpdatedAt = time.Now()
	return uc.playlistRepo.Save(playlist)
}

func (uc *playlistUseCase) EditPlaylistInfo(playlistID uuid.UUID, name, description string) error {
	playlist, err := uc.playlistRepo.FindByID(playlistID)
	if err != nil {
		return fmt.Errorf("playlist not found: %w", err)
	}

	if len(name) > 0 {
		if len(name) < 2 {
			return errors.New("playlist name must be at least 2 characters")
		}
		if len(name) > 100 {
			return errors.New("playlist name is too long")
		}
		playlist.Name = name
	}

	if len(description) > 0 {
		if len(description) > 500 {
			return errors.New("description is too long")
		}
		playlist.Description = description
	}

	playlist.UpdatedAt = time.Now()
	return uc.playlistRepo.Save(playlist)
}

func (uc *playlistUseCase) GetPlaylistTracks(playlistID uuid.UUID) ([]*models.Track, error) {
	if _, err := uc.playlistRepo.FindByID(playlistID); err != nil {
		return nil, fmt.Errorf("playlist not found: %w", err)
	}

	tracks, err := uc.playlistRepo.GetTracks(playlistID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tracks: %w", err)
	}

	return tracks, nil
}

func (uc *playlistUseCase) GetPlaylistWithTracks(playlistID uuid.UUID) (*models.PlaylistTrack, error) {
	playlist, err := uc.playlistRepo.FindByID(playlistID)
	if err != nil {
		return nil, fmt.Errorf("playlist not found: %w", err)
	}

	tracks, err := uc.playlistRepo.GetTracks(playlistID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tracks: %w", err)
	}

	return &models.PlaylistTrack{
		Playlist: *playlist,
		Tracks:   tracks,
	}, nil
}

func (uc *playlistUseCase) GetUserPlaylists(userID uuid.UUID) ([]*models.Playlist, error) {
	if _, err := uc.userRepo.FindByID(userID); err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	playlists, err := uc.playlistRepo.GetUserPlaylists(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user playlists: %w", err)
	}

	return playlists, nil
}

func (uc *playlistUseCase) DeletePlaylist(playlistID, userID uuid.UUID) error {
	playlist, err := uc.playlistRepo.FindByID(playlistID)
	if err != nil {
		return fmt.Errorf("playlist not found: %w", err)
	}

	if playlist.UserID != userID {
		return errors.New("user is not the owner of the playlist")
	}

	if err := uc.playlistRepo.Delete(playlistID); err != nil {
		return fmt.Errorf("failed to delete playlist: %w", err)
	}

	return nil
}
