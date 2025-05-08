package usecases

import (
	"errors"
	"fmt"
	"music-service/internal/models"
	"music-service/internal/repository/interfaces"
	usecaseInterfaces "music-service/internal/usecases/interfaces"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

type albumUseCase struct {
	albumRepo interfaces.AlbumRepository
	trackRepo interfaces.TrackRepository
}

func NewAlbumUseCase(
	albumRepo interfaces.AlbumRepository,
	trackRepo interfaces.TrackRepository,
) usecaseInterfaces.AlbumUseCase {
	return &albumUseCase{
		albumRepo: albumRepo,
		trackRepo: trackRepo,
	}
}

func (uc *albumUseCase) CreateAlbum(title string, releaseDate time.Time, coverURL string) (*models.Album, error) {
	title = strings.TrimSpace(title)
	if utf8.RuneCountInString(title) < 2 {
		return nil, errors.New("album title must be at least 2 characters")
	}
	if utf8.RuneCountInString(title) > 100 {
		return nil, errors.New("album title is too long (max 100 characters)")
	}

	if releaseDate.After(time.Now().Add(24 * time.Hour)) {
		return nil, errors.New("release date cannot be in the future")
	}

	if coverURL != "" {
		if _, err := url.ParseRequestURI(coverURL); err != nil {
			return nil, errors.New("invalid cover URL format")
		}
	}

	existingAlbums, err := uc.albumRepo.ListAll()
	if err != nil {
		return nil, fmt.Errorf("failed to check existing albums: %w", err)
	}

	for _, a := range existingAlbums {
		if strings.EqualFold(a.Title, title) {
			return nil, errors.New("album with this title already exists")
		}
	}

	album := &models.Album{
		Title:       title,
		ReleaseDate: releaseDate,
		CoverURL:    coverURL,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := uc.albumRepo.Save(album); err != nil {
		return nil, fmt.Errorf("failed to save album: %w", err)
	}

	return album, nil
}

func (uc *albumUseCase) AddTrackToAlbum(albumID, trackID uuid.UUID) error {
	album, err := uc.albumRepo.FindByID(albumID)
	if err != nil {
		return fmt.Errorf("album not found: %w", err)
	}

	track, err := uc.trackRepo.FindByID(trackID)
	if err != nil {
		return fmt.Errorf("track not found: %w", err)
	}

	if track.AlbumID != uuid.Nil && track.AlbumID != albumID {
		return errors.New("track already belongs to another album")
	}

	tracks, err := uc.albumRepo.GetTracks(albumID)
	if err != nil {
		return fmt.Errorf("failed to get album tracks: %w", err)
	}

	for _, t := range tracks {
		if t.ID == trackID {
			return errors.New("track already exists in album")
		}
	}

	if len(tracks) >= 50 {
		return errors.New("album cannot contain more than 50 tracks")
	}

	if err := uc.albumRepo.AddTrackToAlbum(albumID, trackID); err != nil {
		return fmt.Errorf("failed to add track to album: %w", err)
	}

	track.AlbumID = albumID
	if err := uc.trackRepo.Save(track); err != nil {
		return fmt.Errorf("failed to update track: %w", err)
	}

	album.UpdatedAt = time.Now()
	return uc.albumRepo.Save(album)
}

func (uc *albumUseCase) RemoveTrackFromAlbum(albumID, trackID uuid.UUID) error {
	album, err := uc.albumRepo.FindByID(albumID)
	if err != nil {
		return fmt.Errorf("album not found: %w", err)
	}

	track, err := uc.trackRepo.FindByID(trackID)
	if err != nil {
		return fmt.Errorf("track not found: %w", err)
	}

	if track.AlbumID != albumID {
		return errors.New("track does not belong to this album")
	}

	if err := uc.albumRepo.RemoveTrackFromAlbum(albumID, trackID); err != nil {
		return fmt.Errorf("failed to remove track from album: %w", err)
	}

	track.AlbumID = uuid.Nil
	if err := uc.trackRepo.Save(track); err != nil {
		return fmt.Errorf("failed to update track: %w", err)
	}

	album.UpdatedAt = time.Now()
	return uc.albumRepo.Save(album)
}

func (uc *albumUseCase) GetAlbumDetails(albumID uuid.UUID) (*models.Album, []*models.Track, error) {
	album, err := uc.albumRepo.FindByID(albumID)
	if err != nil {
		return nil, nil, fmt.Errorf("album not found: %w", err)
	}

	tracks, err := uc.albumRepo.GetTracks(albumID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get album tracks: %w", err)
	}

	sortTracks(tracks)

	return album, tracks, nil
}

func (uc *albumUseCase) UpdateAlbumInfo(albumID uuid.UUID, title, coverURL string, releaseDate time.Time) error {
	album, err := uc.albumRepo.FindByID(albumID)
	if err != nil {
		return fmt.Errorf("album not found: %w", err)
	}

	if title != "" {
		title = strings.TrimSpace(title)
		if utf8.RuneCountInString(title) < 2 {
			return errors.New("album title must be at least 2 characters")
		}
		if utf8.RuneCountInString(title) > 100 {
			return errors.New("album title is too long (max 100 characters)")
		}
		album.Title = title
	}

	if coverURL != "" {
		if _, err := url.ParseRequestURI(coverURL); err != nil {
			return errors.New("invalid cover URL format")
		}
		album.CoverURL = coverURL
	}

	if !releaseDate.IsZero() {
		if releaseDate.After(time.Now().Add(24 * time.Hour)) {
			return errors.New("release date cannot be in the future")
		}
		album.ReleaseDate = releaseDate
	}

	album.UpdatedAt = time.Now()
	return uc.albumRepo.Save(album)
}

func sortTracks(tracks []*models.Track) {
	for i := 0; i < len(tracks)-1; i++ {
		for j := i + 1; j < len(tracks); j++ {
			if tracks[i].ID.String() > tracks[j].ID.String() {
				tracks[i], tracks[j] = tracks[j], tracks[i]
			}
		}
	}
}
