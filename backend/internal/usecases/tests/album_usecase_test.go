package usecases_test

import (
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"music-service/internal/models"
	"music-service/internal/repository/interfaces/mocks"
	"music-service/internal/usecases"
)

func TestCreateAlbum(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	albumRepo := mocks.NewMockAlbumRepository(ctrl)
	trackRepo := mocks.NewMockTrackRepository(ctrl)

	uc := usecases.NewAlbumUseCase(albumRepo, trackRepo)

	tests := []struct {
		name        string
		title       string
		releaseDate time.Time
		coverURL    string
		mockSetup   func()
		expected    *models.Album
		expectedErr string
	}{
		{
			name:        "successful creation",
			title:       "Valid Album",
			releaseDate: time.Now().Add(-24 * time.Hour),
			coverURL:    "https://example.com/cover.jpg",
			mockSetup: func() {
				albumRepo.EXPECT().ListAll().Return([]*models.Album{}, nil)
				albumRepo.EXPECT().Save(gomock.Any()).DoAndReturn(func(album *models.Album) error {
					album.ID = uuid.New()
					return nil
				})
			},
			expected: &models.Album{
				ID:          uuid.New(),
				Title:       "Valid Album",
				ReleaseDate: time.Now().Add(-24 * time.Hour),
				CoverURL:    "https://example.com/cover.jpg",
			},
		},
		{
			name:        "title too short",
			title:       "A",
			releaseDate: time.Now(),
			expectedErr: "album title must be at least 2 characters",
		},
		{
			name:        "invalid cover URL",
			title:       "Valid Album",
			releaseDate: time.Now(),
			coverURL:    "invalid-url",
			expectedErr: "invalid cover URL format",
		},
		{
			name:        "duplicate title",
			title:       "Existing Album",
			releaseDate: time.Now(),
			mockSetup: func() {
				albumRepo.EXPECT().ListAll().Return([]*models.Album{
					{Title: "Existing Album"},
				}, nil)
			},
			expectedErr: "album with this title already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			result, err := uc.CreateAlbum(tt.title, tt.releaseDate, tt.coverURL)

			if tt.expectedErr != "" {
				assert.ErrorContains(t, err, tt.expectedErr)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expected.Title, result.Title)
				assert.Equal(t, tt.expected.CoverURL, result.CoverURL)
				assert.WithinDuration(t, tt.expected.ReleaseDate, result.ReleaseDate, time.Second)
				assert.NotZero(t, result.ID)
				assert.NotZero(t, result.CreatedAt)
				assert.NotZero(t, result.UpdatedAt)
			}
		})
	}
}

func TestAddTrackToAlbum(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	albumRepo := mocks.NewMockAlbumRepository(ctrl)
	trackRepo := mocks.NewMockTrackRepository(ctrl)
	uc := usecases.NewAlbumUseCase(albumRepo, trackRepo)

	albumID := uuid.New()
	trackID := uuid.New()

	tests := []struct {
		name        string
		albumID     uuid.UUID
		trackID     uuid.UUID
		mockSetup   func()
		expectedErr string
	}{
		{
			name:    "successful add",
			albumID: albumID,
			trackID: trackID,
			mockSetup: func() {
				albumRepo.EXPECT().FindByID(albumID).Return(&models.Album{ID: albumID}, nil)
				trackRepo.EXPECT().FindByID(trackID).Return(&models.Track{ID: trackID}, nil)
				albumRepo.EXPECT().GetTracks(albumID).Return([]*models.Track{}, nil)
				albumRepo.EXPECT().AddTrackToAlbum(albumID, trackID).Return(nil)
				trackRepo.EXPECT().Save(gomock.Any()).Return(nil)
				albumRepo.EXPECT().Save(gomock.Any()).Return(nil)
			},
		},
		{
			name:    "album not found",
			albumID: uuid.New(),
			trackID: trackID,
			mockSetup: func() {
				albumRepo.EXPECT().FindByID(gomock.Any()).Return(nil, errors.New("not found"))
			},
			expectedErr: "album not found",
		},
		{
			name:    "track already in another album",
			albumID: albumID,
			trackID: trackID,
			mockSetup: func() {
				albumRepo.EXPECT().FindByID(albumID).Return(&models.Album{ID: albumID}, nil)
				trackRepo.EXPECT().FindByID(trackID).Return(&models.Track{ID: trackID, AlbumID: uuid.New()}, nil)
			},
			expectedErr: "track already belongs to another album",
		},
		{
			name:    "album full (50 tracks)",
			albumID: albumID,
			trackID: trackID,
			mockSetup: func() {
				albumRepo.EXPECT().FindByID(albumID).Return(&models.Album{ID: albumID}, nil)
				trackRepo.EXPECT().FindByID(trackID).Return(&models.Track{ID: trackID}, nil)

				// Создаем 50 уникальных треков
				tracks := make([]*models.Track, 50)
				for i := range tracks {
					tracks[i] = &models.Track{ID: uuid.New()} // ID, отличные от trackID
				}
				albumRepo.EXPECT().GetTracks(albumID).Return(tracks, nil)
			},
			expectedErr: "album cannot contain more than 50 tracks",
		},
		{
			name:    "track already exists in album",
			albumID: albumID,
			trackID: trackID,
			mockSetup: func() {
				albumRepo.EXPECT().FindByID(albumID).Return(&models.Album{ID: albumID}, nil)
				trackRepo.EXPECT().FindByID(trackID).Return(&models.Track{ID: trackID}, nil)
				albumRepo.EXPECT().GetTracks(albumID).Return([]*models.Track{
					{ID: trackID},
				}, nil)
			},
			expectedErr: "track already exists in album",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			err := uc.AddTrackToAlbum(tt.albumID, tt.trackID)

			if tt.expectedErr != "" {
				assert.ErrorContains(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRemoveTrackFromAlbum(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	albumRepo := mocks.NewMockAlbumRepository(ctrl)
	trackRepo := mocks.NewMockTrackRepository(ctrl)
	uc := usecases.NewAlbumUseCase(albumRepo, trackRepo)

	albumID := uuid.New()
	trackID := uuid.New()

	tests := []struct {
		name        string
		albumID     uuid.UUID
		trackID     uuid.UUID
		mockSetup   func()
		expectedErr string
	}{
		{
			name:    "successful remove",
			albumID: albumID,
			trackID: trackID,
			mockSetup: func() {
				albumRepo.EXPECT().FindByID(albumID).Return(&models.Album{ID: albumID}, nil)
				trackRepo.EXPECT().FindByID(trackID).Return(&models.Track{ID: trackID, AlbumID: albumID}, nil)
				albumRepo.EXPECT().RemoveTrackFromAlbum(albumID, trackID).Return(nil)
				trackRepo.EXPECT().Save(gomock.Any()).Return(nil)
				albumRepo.EXPECT().Save(gomock.Any()).Return(nil)
			},
		},
		{
			name:    "track not in album",
			albumID: albumID,
			trackID: trackID,
			mockSetup: func() {
				albumRepo.EXPECT().FindByID(albumID).Return(&models.Album{ID: albumID}, nil)
				trackRepo.EXPECT().FindByID(trackID).Return(&models.Track{ID: trackID, AlbumID: uuid.New()}, nil)
			},
			expectedErr: "track does not belong to this album",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			err := uc.RemoveTrackFromAlbum(tt.albumID, tt.trackID)

			if tt.expectedErr != "" {
				assert.ErrorContains(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetAlbumDetails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	albumRepo := mocks.NewMockAlbumRepository(ctrl)
	trackRepo := mocks.NewMockTrackRepository(ctrl)
	uc := usecases.NewAlbumUseCase(albumRepo, trackRepo)

	albumID := uuid.New()
	track := &models.Track{ID: uuid.New(), Title: "Track 1"}

	tests := []struct {
		name        string
		albumID     uuid.UUID
		mockSetup   func()
		expected    *models.Album
		expectedErr string
	}{
		{
			name:    "successful get",
			albumID: albumID,
			mockSetup: func() {
				albumRepo.EXPECT().FindByID(albumID).Return(&models.Album{
					ID:    albumID,
					Title: "Test Album",
				}, nil)
				albumRepo.EXPECT().GetTracks(albumID).Return([]*models.Track{track}, nil)
			},
			expected: &models.Album{
				ID:    albumID,
				Title: "Test Album",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			album, tracks, err := uc.GetAlbumDetails(tt.albumID)

			if tt.expectedErr != "" {
				assert.ErrorContains(t, err, tt.expectedErr)
				assert.Nil(t, album)
				assert.Nil(t, tracks)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.ID, album.ID)
				assert.Equal(t, tt.expected.Title, album.Title)
				assert.Len(t, tracks, 1)
				assert.Equal(t, track.ID, tracks[0].ID)
			}
		})
	}
}

func TestUpdateAlbumInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	albumRepo := mocks.NewMockAlbumRepository(ctrl)
	trackRepo := mocks.NewMockTrackRepository(ctrl)
	uc := usecases.NewAlbumUseCase(albumRepo, trackRepo)

	albumID := uuid.New()
	validURL := "https://example.com/new-cover.jpg"

	tests := []struct {
		name        string
		albumID     uuid.UUID
		title       string
		coverURL    string
		releaseDate time.Time
		mockSetup   func()
		expectedErr string
	}{
		{
			name:     "successful update",
			albumID:  albumID,
			title:    "New Title",
			coverURL: validURL,
			mockSetup: func() {
				albumRepo.EXPECT().FindByID(albumID).Return(&models.Album{
					ID:    albumID,
					Title: "Old Title",
				}, nil)
				albumRepo.EXPECT().Save(gomock.Any()).Return(nil)
			},
		},
		{
			name:    "invalid title",
			albumID: albumID,
			title:   "A",
			mockSetup: func() {
				albumRepo.EXPECT().FindByID(albumID).Return(&models.Album{
					ID:    albumID,
					Title: "Old Title",
				}, nil)
			},
			expectedErr: "album title must be at least 2 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			err := uc.UpdateAlbumInfo(tt.albumID, tt.title, tt.coverURL, tt.releaseDate)

			if tt.expectedErr != "" {
				assert.ErrorContains(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
