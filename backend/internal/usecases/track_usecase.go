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

type trackUseCase struct {
	trackRepo   interfaces.TrackRepository
	historyRepo interfaces.HistoryRepository
	albumRepo   interfaces.AlbumRepository
}

func NewTrackUseCase(
	trackRepo interfaces.TrackRepository,
	historyRepo interfaces.HistoryRepository,
	albumRepo interfaces.AlbumRepository,
) usecaseInterfaces.TrackUseCase {
	return &trackUseCase{
		trackRepo:   trackRepo,
		historyRepo: historyRepo,
		albumRepo:   albumRepo,
	}
}

func (uc *trackUseCase) SearchTracks(query string) ([]*models.Track, error) {
	if len(query) < 3 {
		return nil, errors.New("search query must be at least 2 characters")
	}

	tracks, err := uc.trackRepo.Search(query)
	if err != nil {
		return nil, fmt.Errorf("failed to search tracks: %w", err)
	}

	return tracks, nil
}

func (uc *trackUseCase) PlayTrack(userID, trackID uuid.UUID) error {
	track, err := uc.trackRepo.FindByID(trackID)
	if err != nil {
		return fmt.Errorf("track not found: %w", err)
	}
	if track.AlbumID != uuid.Nil {
		if _, err := uc.albumRepo.FindByID(track.AlbumID); err != nil {
			return fmt.Errorf("album not found: %w", err)
		}
	}
	if err := uc.historyRepo.AddEntry(userID, trackID); err != nil {
		return fmt.Errorf("failed to record playback: %w", err)
	}

	go func() {
		_ = uc.trackRepo.IncrementPlayCount(trackID)
	}()

	return nil
}

func (uc *trackUseCase) GetTrackDetails(trackID uuid.UUID) (*models.Track, error) {
	track, err := uc.trackRepo.FindByID(trackID)
	if err != nil {
		return nil, fmt.Errorf("track not found: %w", err)
	}
	return track, nil
}

func (uc *trackUseCase) UpdateTrackMetadata(trackID uuid.UUID, metadata map[string]interface{}) error {
	track, err := uc.trackRepo.FindByID(trackID)
	if err != nil {
		return fmt.Errorf("track not found: %w", err)
	}

	if title, ok := metadata["title"].(string); ok && title != "" {
		if len(title) > 100 {
			return errors.New("title is too long")
		}
		track.Title = title
	}

	if artist, ok := metadata["artist_name"].(string); ok && artist != "" {
		track.ArtistName = artist
	}

	if albumID, ok := metadata["album_id"].(uuid.UUID); ok && albumID != uuid.Nil {
		if _, err := uc.albumRepo.FindByID(albumID); err != nil {
			return fmt.Errorf("album not found: %w", err)
		}
		track.AlbumID = albumID
	}

	track.UpdatedAt = time.Now()
	return uc.trackRepo.Save(track)
}

func (uc *trackUseCase) DeleteTrack(trackID uuid.UUID) error {
	if _, err := uc.trackRepo.FindByID(trackID); err != nil {
		return fmt.Errorf("track not found: %w", err)
	}

	return uc.trackRepo.Delete(trackID)
}

func formatDuration(seconds int) string {
	minutes := seconds / 60
	remainingSeconds := seconds % 60
	return fmt.Sprintf("%d:%02d", minutes, remainingSeconds)
}
