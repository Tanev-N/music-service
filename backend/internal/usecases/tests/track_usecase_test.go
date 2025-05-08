package usecases_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"music-service/internal/models"
	"music-service/internal/repository/interfaces/mocks"
	"music-service/internal/usecases"
)

func TestSearchTracks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	trackRepo := mocks.NewMockTrackRepository(ctrl)
	historyRepo := mocks.NewMockHistoryRepository(ctrl)
	albumRepo := mocks.NewMockAlbumRepository(ctrl)
	uc := usecases.NewTrackUseCase(trackRepo, historyRepo, albumRepo)

	// Задаем фиксированный UUID для ожидаемого трека
	expectedTrackID := uuid.MustParse("1e872412-8a31-4790-b118-ce335a3f4e84")

	tests := []struct {
		name        string
		query       string
		mockSetup   func()
		expected    []*models.Track
		expectedErr string
	}{
		{
			name:  "successful search",
			query: "test",
			mockSetup: func() {
				trackRepo.EXPECT().Search("test").Return([]*models.Track{
					{ID: expectedTrackID, Title: "Test Track"}, // Используем фиксированный UUID
				}, nil)
			},
			expected: []*models.Track{
				{ID: expectedTrackID, Title: "Test Track"}, // Используем фиксированный UUID
			},
		},
		{
			name:        "query too short",
			query:       "te",
			expectedErr: "search query must be at least 2 characters",
		},
		{
			name:  "search error",
			query: "error",
			mockSetup: func() {
				trackRepo.EXPECT().Search("error").Return(nil, errors.New("search failed"))
			},
			expectedErr: "failed to search tracks",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			result, err := uc.SearchTracks(tt.query)

			if tt.expectedErr != "" {
				assert.ErrorContains(t, err, tt.expectedErr)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestPlayTrack(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	trackRepo := mocks.NewMockTrackRepository(ctrl)
	historyRepo := mocks.NewMockHistoryRepository(ctrl)
	albumRepo := mocks.NewMockAlbumRepository(ctrl)
	uc := usecases.NewTrackUseCase(trackRepo, historyRepo, albumRepo)

	// Фиксированный UUID для пользователя и трека
	userID := uuid.MustParse("1f22acbc-d723-4f92-a06f-8a583b4d45ea")
	trackID := uuid.MustParse("94db4f59-1757-4e9d-bf14-c6df52eafab4")

	tests := []struct {
		name        string
		userID      uuid.UUID
		trackID     uuid.UUID
		mockSetup   func()
		expectedErr string
	}{
		{
			name:    "successful play",
			userID:  userID,
			trackID: trackID,
			mockSetup: func() {
				trackRepo.EXPECT().FindByID(trackID).Return(&models.Track{ID: trackID}, nil)
				historyRepo.EXPECT().AddEntry(userID, trackID).Return(nil)
				trackRepo.EXPECT().IncrementPlayCount(trackID).Return(nil) // Ожидаем вызов IncrementPlayCount
			},
		},
		{
			name:    "track not found",
			userID:  userID,
			trackID: trackID,
			mockSetup: func() {
				trackRepo.EXPECT().FindByID(trackID).Return(nil, errors.New("not found"))
			},
			expectedErr: "track not found",
		},
		{
			name:    "history record failed",
			userID:  userID,
			trackID: trackID,
			mockSetup: func() {
				trackRepo.EXPECT().FindByID(trackID).Return(&models.Track{ID: trackID}, nil)
				historyRepo.EXPECT().AddEntry(userID, trackID).Return(errors.New("failed"))
			},
			expectedErr: "failed to record playback",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			err := uc.PlayTrack(tt.userID, tt.trackID)

			if tt.expectedErr != "" {
				assert.ErrorContains(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetTrackDetails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	trackRepo := mocks.NewMockTrackRepository(ctrl)
	historyRepo := mocks.NewMockHistoryRepository(ctrl)
	albumRepo := mocks.NewMockAlbumRepository(ctrl)
	uc := usecases.NewTrackUseCase(trackRepo, historyRepo, albumRepo)

	// Фиксированный UUID для трека
	trackID := uuid.MustParse("1e872412-8a31-4790-b118-ce335a3f4e84")

	tests := []struct {
		name        string
		trackID     uuid.UUID
		mockSetup   func()
		expected    *models.Track
		expectedErr string
	}{
		{
			name:    "successful get",
			trackID: trackID,
			mockSetup: func() {
				trackRepo.EXPECT().FindByID(trackID).Return(&models.Track{
					ID:    trackID,
					Title: "Test Track",
				}, nil)
			},
			expected: &models.Track{
				ID:    trackID,
				Title: "Test Track",
			},
		},
		{
			name:    "track not found",
			trackID: uuid.New(),
			mockSetup: func() {
				trackRepo.EXPECT().FindByID(gomock.Any()).Return(nil, errors.New("not found"))
			},
			expectedErr: "track not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			result, err := uc.GetTrackDetails(tt.trackID)

			if tt.expectedErr != "" {
				assert.ErrorContains(t, err, tt.expectedErr)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestUpdateTrackMetadata(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	trackRepo := mocks.NewMockTrackRepository(ctrl)
	historyRepo := mocks.NewMockHistoryRepository(ctrl)
	albumRepo := mocks.NewMockAlbumRepository(ctrl)
	uc := usecases.NewTrackUseCase(trackRepo, historyRepo, albumRepo)

	trackID := uuid.MustParse("1e872412-8a31-4790-b118-ce335a3f4e84")
	albumId := uuid.New()

	tests := []struct {
		name        string
		trackID     uuid.UUID
		metadata    map[string]interface{}
		mockSetup   func()
		expectedErr string
	}{
		{
			name:    "title too long",
			trackID: trackID,
			metadata: map[string]interface{}{
				"title": "This title is way too long and exceeds the maximum allowed length of 100 characters which should trigger an error",
			},
			mockSetup: func() {
				trackRepo.EXPECT().FindByID(trackID).Return(&models.Track{ID: trackID}, nil)
			},
			expectedErr: "title is too long",
		},
		{
			name:    "album not found",
			trackID: trackID,
			metadata: map[string]interface{}{
				"album_id": albumId,
			},
			mockSetup: func() {
				trackRepo.EXPECT().FindByID(trackID).Return(&models.Track{ID: trackID}, nil)
				albumRepo.EXPECT().FindByID(albumId).Return(nil, errors.New("not found"))
			},
			expectedErr: "album not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			err := uc.UpdateTrackMetadata(tt.trackID, tt.metadata)

			if tt.expectedErr != "" {
				assert.ErrorContains(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteTrack(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	trackRepo := mocks.NewMockTrackRepository(ctrl)
	historyRepo := mocks.NewMockHistoryRepository(ctrl)
	albumRepo := mocks.NewMockAlbumRepository(ctrl)
	uc := usecases.NewTrackUseCase(trackRepo, historyRepo, albumRepo)

	trackID := uuid.MustParse("1e872412-8a31-4790-b118-ce335a3f4e84")

	tests := []struct {
		name        string
		trackID     uuid.UUID
		mockSetup   func()
		expectedErr string
	}{
		{
			name:    "successful delete",
			trackID: trackID,
			mockSetup: func() {
				trackRepo.EXPECT().FindByID(trackID).Return(&models.Track{ID: trackID}, nil)
				trackRepo.EXPECT().Delete(trackID).Return(nil)
			},
		},
		{
			name:    "track not found",
			trackID: uuid.New(),
			mockSetup: func() {
				trackRepo.EXPECT().FindByID(gomock.Any()).Return(nil, errors.New("not found"))
			},
			expectedErr: "track not found",
		},
		{
			name:    "delete failed",
			trackID: trackID,
			mockSetup: func() {
				trackRepo.EXPECT().FindByID(trackID).Return(&models.Track{ID: trackID}, nil)
				trackRepo.EXPECT().Delete(trackID).Return(errors.New("delete failed"))
			},
			expectedErr: "delete failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			err := uc.DeleteTrack(tt.trackID)

			if tt.expectedErr != "" {
				assert.ErrorContains(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
