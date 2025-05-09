package usecases

import (
	"errors"
	"fmt"
	"io"
	"log"
	"music-service/internal/models"
	"music-service/internal/repository/interfaces"
	usecaseInterfaces "music-service/internal/usecases/interfaces"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type trackUseCase struct {
	trackRepo     interfaces.TrackRepository
	historyRepo   interfaces.HistoryRepository
	albumRepo     interfaces.AlbumRepository
	maxFileSizeMB int
	allowedTypes  []string
}

func NewTrackUseCase(
	trackRepo interfaces.TrackRepository,
	historyRepo interfaces.HistoryRepository,
	albumRepo interfaces.AlbumRepository,
	maxFileSizeMB int,
	allowedTypes []string,
) usecaseInterfaces.TrackUseCase {
	return &trackUseCase{
		trackRepo:     trackRepo,
		historyRepo:   historyRepo,
		albumRepo:     albumRepo,
		maxFileSizeMB: maxFileSizeMB,
		allowedTypes:  allowedTypes,
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

func (uc *trackUseCase) GetTrackDetails(id uuid.UUID) (*models.TrackDetails, error) {
	track, err := uc.trackRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("track not found: %w", err)
	}

	var album *models.Album
	if track.AlbumID != uuid.Nil {
		album, err = uc.albumRepo.FindByID(track.AlbumID)
		if err != nil {
			log.Printf("could not find album for track %s: %v", id, err)
		}
	}

	genres, err := uc.trackRepo.GetGenresForTrack(id)
	if err != nil {
		log.Printf("could not find genres for track %s: %v", id, err)
	}

	playCount, err := uc.historyRepo.GetPlayCount(id)
	if err != nil {
		log.Printf("could not get play count for track %s: %v", id, err)
	}

	return &models.TrackDetails{
		ID:         track.ID,
		Title:      track.Title,
		ArtistName: track.ArtistName,
		Duration:   track.Duration,
		FilePath:   track.FilePath,
		MimeType:   "audio/mpeg",
		CoverURL:   track.CoverURL,
		AddedDate:  track.AddedDate,
		CreatedAt:  track.AddedDate,
		UpdatedAt:  track.UpdatedAt,
		PlayCount:  playCount,
		Album:      album,
		Genres:     genres,
	}, nil
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

func (uc *trackUseCase) UploadTrack(fileReader io.Reader, fileSize int64, metadata models.TrackUploadMetadata) (*models.Track, error) {
	maxSizeBytes := int64(uc.maxFileSizeMB * 1024 * 1024)
	if fileSize > maxSizeBytes {
		return nil, fmt.Errorf("размер файла превышает максимально допустимый (%d МБ)", uc.maxFileSizeMB)
	}

	if metadata.Title == "" {
		return nil, errors.New("название трека не может быть пустым")
	}

	if metadata.ArtistName == "" {
		return nil, errors.New("имя исполнителя не может быть пустым")
	}

	if metadata.AlbumID == uuid.Nil {
		return nil, errors.New("необходимо указать альбом для трека")
	}

	if _, err := uc.albumRepo.FindByID(metadata.AlbumID); err != nil {
		return nil, fmt.Errorf("альбом не найден: %w", err)
	}

	trackID := uuid.New()

	filePath, err := uc.trackRepo.SaveTrackFile(trackID, fileReader, fileSize)
	if err != nil {
		return nil, fmt.Errorf("ошибка при сохранении файла: %w", err)
	}

	now := time.Now()
	track := &models.Track{
		ID:         trackID,
		Title:      metadata.Title,
		Duration:   metadata.Duration,
		FilePath:   filePath,
		AlbumID:    metadata.AlbumID,
		ArtistName: metadata.ArtistName,
		CoverURL:   metadata.CoverURL,
		AddedDate:  now,
		UpdatedAt:  now,
		PlayCount:  0,
	}

	if err := uc.trackRepo.Save(track); err != nil {
		return nil, fmt.Errorf("ошибка при сохранении метаданных трека: %w", err)
	}

	return track, nil
}

func (uc *trackUseCase) GetTrackFilePath(trackID uuid.UUID) (string, error) {
	track, err := uc.trackRepo.FindByID(trackID)
	if err != nil {
		return "", fmt.Errorf("трек не найден: %w", err)
	}

	tracksDir := uc.trackRepo.GetStorageDir()
	if tracksDir == "" {
		return "", errors.New("путь к директории треков не настроен")
	}

	fullPath := filepath.Join(tracksDir, track.FilePath)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("файл трека не найден: %w", err)
	}

	return fullPath, nil
}
