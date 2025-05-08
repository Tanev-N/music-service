package usecases_test

import (
	"errors"
	"music-service/internal/models"
	"music-service/internal/repository/interfaces/mocks"
	"testing"
	"time"

	"music-service/internal/usecases"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAlbumUseCase_CreateAlbum(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAlbumRepo := mocks.NewMockAlbumRepository(ctrl)
	mockTrackRepo := mocks.NewMockTrackRepository(ctrl)

	useCase := usecases.NewAlbumUseCase(mockAlbumRepo, mockTrackRepo)

	title := "Test Album"
	releaseDate := time.Now().Add(-24 * time.Hour) // вчера
	coverURL := "http://example.com/cover.jpg"

	t.Run("album_already_exists", func(t *testing.T) {
		existingAlbums := []*models.Album{
			{
				ID:          uuid.New(),
				Title:       "Test Album",
				ReleaseDate: time.Now(),
				CoverURL:    "http://example.com/another_cover.jpg",
			},
		}

		mockAlbumRepo.EXPECT().
			ListAll().
			Return(existingAlbums, nil)

		album, err := useCase.CreateAlbum(title, releaseDate, coverURL)
		assert.Error(t, err)
		assert.Nil(t, album)
		assert.Contains(t, err.Error(), "already exists")
	})

	// Проверка успешного создания
	t.Run("success", func(t *testing.T) {
		mockAlbumRepo.EXPECT().
			ListAll().
			Return([]*models.Album{}, nil)

		mockAlbumRepo.EXPECT().
			Save(gomock.Any()).
			DoAndReturn(func(album *models.Album) error {
				assert.Equal(t, title, album.Title)
				assert.Equal(t, releaseDate, album.ReleaseDate)
				assert.Equal(t, coverURL, album.CoverURL)
				return nil
			})

		album, err := useCase.CreateAlbum(title, releaseDate, coverURL)
		assert.NoError(t, err)
		assert.NotNil(t, album)
		assert.Equal(t, title, album.Title)
		assert.Equal(t, releaseDate, album.ReleaseDate)
		assert.Equal(t, coverURL, album.CoverURL)
	})

	// Ошибка при сохранении альбома
	t.Run("save_error", func(t *testing.T) {
		mockAlbumRepo.EXPECT().
			ListAll().
			Return([]*models.Album{}, nil)

		mockAlbumRepo.EXPECT().
			Save(gomock.Any()).
			Return(errors.New("db error"))

		album, err := useCase.CreateAlbum(title, releaseDate, coverURL)
		assert.Error(t, err)
		assert.Nil(t, album)
		assert.Contains(t, err.Error(), "failed to save album")
	})

	// Проверка валидации: название слишком короткое
	t.Run("title_too_short", func(t *testing.T) {
		album, err := useCase.CreateAlbum("A", releaseDate, coverURL)
		assert.Error(t, err)
		assert.Nil(t, album)
		assert.Contains(t, err.Error(), "at least 2 characters")
	})

	// Проверка валидации: дата выпуска в будущем
	t.Run("future_release_date", func(t *testing.T) {
		futureDate := time.Now().Add(48 * time.Hour) // через 2 дня
		album, err := useCase.CreateAlbum(title, futureDate, coverURL)
		assert.Error(t, err)
		assert.Nil(t, album)
		assert.Contains(t, err.Error(), "release date cannot be in the future")
	})

	// Проверка валидации: неверный формат URL обложки
	t.Run("invalid_cover_url", func(t *testing.T) {
		album, err := useCase.CreateAlbum(title, releaseDate, "invalid-url")
		assert.Error(t, err)
		assert.Nil(t, album)
		assert.Contains(t, err.Error(), "invalid cover URL format")
	})
}

func TestAlbumUseCase_AddTrackToAlbum(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAlbumRepo := mocks.NewMockAlbumRepository(ctrl)
	mockTrackRepo := mocks.NewMockTrackRepository(ctrl)

	useCase := usecases.NewAlbumUseCase(mockAlbumRepo, mockTrackRepo)

	albumID := uuid.New()
	trackID := uuid.New()
	anotherAlbumID := uuid.New()

	album := &models.Album{
		ID:          albumID,
		Title:       "Test Album",
		ReleaseDate: time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Успешное добавление трека
	t.Run("success", func(t *testing.T) {
		track := &models.Track{
			ID:        trackID,
			Title:     "Test Track",
			AlbumID:   uuid.Nil, // трек не принадлежит ни к одному альбому
			UpdatedAt: time.Now(),
			PlayCount: 0,
		}

		mockAlbumRepo.EXPECT().
			FindByID(albumID).
			Return(album, nil)

		mockTrackRepo.EXPECT().
			FindByID(trackID).
			Return(track, nil)

		mockAlbumRepo.EXPECT().
			GetTracks(albumID).
			Return([]*models.Track{}, nil)

		mockAlbumRepo.EXPECT().
			AddTrackToAlbum(albumID, trackID).
			Return(nil)

		mockTrackRepo.EXPECT().
			Save(gomock.Any()).
			DoAndReturn(func(track *models.Track) error {
				assert.Equal(t, albumID, track.AlbumID)
				return nil
			})

		mockAlbumRepo.EXPECT().
			Save(gomock.Any()).
			DoAndReturn(func(a *models.Album) error {
				assert.Equal(t, albumID, a.ID)
				return nil
			})

		err := useCase.AddTrackToAlbum(albumID, trackID)
		assert.NoError(t, err)
	})

	// Трек уже принадлежит другому альбому
	t.Run("track_in_another_album", func(t *testing.T) {
		track := &models.Track{
			ID:        trackID,
			Title:     "Test Track",
			AlbumID:   anotherAlbumID, // трек уже в другом альбоме
			UpdatedAt: time.Now(),
		}

		mockAlbumRepo.EXPECT().
			FindByID(albumID).
			Return(album, nil)

		mockTrackRepo.EXPECT().
			FindByID(trackID).
			Return(track, nil)

		err := useCase.AddTrackToAlbum(albumID, trackID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already belongs to another album")
	})

	// Трек уже есть в альбоме
	t.Run("track_already_in_album", func(t *testing.T) {
		track := &models.Track{
			ID:        trackID,
			Title:     "Test Track",
			AlbumID:   albumID, // трек уже в этом альбоме
			UpdatedAt: time.Now(),
		}

		mockAlbumRepo.EXPECT().
			FindByID(albumID).
			Return(album, nil)

		mockTrackRepo.EXPECT().
			FindByID(trackID).
			Return(track, nil)

		mockAlbumRepo.EXPECT().
			GetTracks(albumID).
			Return([]*models.Track{track}, nil)

		err := useCase.AddTrackToAlbum(albumID, trackID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already exists in album")
	})

	// Альбом не найден
	t.Run("album_not_found", func(t *testing.T) {
		mockAlbumRepo.EXPECT().
			FindByID(albumID).
			Return(nil, errors.New("not found"))

		err := useCase.AddTrackToAlbum(albumID, trackID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "album not found")
	})

	// Трек не найден
	t.Run("track_not_found", func(t *testing.T) {
		mockAlbumRepo.EXPECT().
			FindByID(albumID).
			Return(album, nil)

		mockTrackRepo.EXPECT().
			FindByID(trackID).
			Return(nil, errors.New("not found"))

		err := useCase.AddTrackToAlbum(albumID, trackID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "track not found")
	})
}

func TestAlbumUseCase_GetAlbumDetails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAlbumRepo := mocks.NewMockAlbumRepository(ctrl)
	mockTrackRepo := mocks.NewMockTrackRepository(ctrl)

	useCase := usecases.NewAlbumUseCase(mockAlbumRepo, mockTrackRepo)

	albumID := uuid.New()
	album := &models.Album{
		ID:          albumID,
		Title:       "Test Album",
		ReleaseDate: time.Now(),
	}

	tracks := []*models.Track{
		{
			ID:        uuid.New(),
			Title:     "Track 1",
			AlbumID:   albumID,
			Duration:  180,
			PlayCount: 10,
		},
		{
			ID:        uuid.New(),
			Title:     "Track 2",
			AlbumID:   albumID,
			Duration:  240,
			PlayCount: 5,
		},
	}

	t.Run("success", func(t *testing.T) {
		mockAlbumRepo.EXPECT().
			FindByID(albumID).
			Return(album, nil)

		mockAlbumRepo.EXPECT().
			GetTracks(albumID).
			Return(tracks, nil)

		foundAlbum, foundTracks, err := useCase.GetAlbumDetails(albumID)
		assert.NoError(t, err)
		assert.Equal(t, album, foundAlbum)
		assert.Equal(t, len(tracks), len(foundTracks))
	})

	t.Run("album_not_found", func(t *testing.T) {
		mockAlbumRepo.EXPECT().
			FindByID(albumID).
			Return(nil, errors.New("not found"))

		foundAlbum, foundTracks, err := useCase.GetAlbumDetails(albumID)
		assert.Error(t, err)
		assert.Nil(t, foundAlbum)
		assert.Nil(t, foundTracks)
		assert.Contains(t, err.Error(), "album not found")
	})

	t.Run("tracks_error", func(t *testing.T) {
		mockAlbumRepo.EXPECT().
			FindByID(albumID).
			Return(album, nil)

		mockAlbumRepo.EXPECT().
			GetTracks(albumID).
			Return(nil, errors.New("db error"))

		foundAlbum, foundTracks, err := useCase.GetAlbumDetails(albumID)
		assert.Error(t, err)
		assert.Nil(t, foundAlbum)
		assert.Nil(t, foundTracks)
		assert.Contains(t, err.Error(), "failed to get album tracks")
	})
}
