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

func TestRecordPlayback(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	historyRepo := mocks.NewMockHistoryRepository(ctrl)
	trackRepo := mocks.NewMockTrackRepository(ctrl)
	uc := usecases.NewHistoryUseCase(historyRepo, trackRepo)

	userID := uuid.MustParse("1f22acbc-d723-4f92-a06f-8a583b4d45ea")
	trackID := uuid.MustParse("2f46e5ce-8b10-4994-bc10-d33637e98ebf")
	now := time.Now()

	tests := []struct {
		name        string
		userID      uuid.UUID
		trackID     uuid.UUID
		mockSetup   func()
		expectedErr string
	}{
		{
			name:    "successful recording",
			userID:  userID,
			trackID: trackID,
			mockSetup: func() {
				trackRepo.EXPECT().FindByID(trackID).Return(&models.Track{
					ID:       trackID,
					Duration: 180, // 3 минуты
				}, nil)
				historyRepo.EXPECT().GetHistory(userID).Return([]*models.ListeningHistory{}, nil)
				historyRepo.EXPECT().AddEntry(userID, trackID).Return(nil)
			},
		},
		{
			name:    "track not found",
			userID:  userID,
			trackID: uuid.New(),
			mockSetup: func() {
				trackRepo.EXPECT().FindByID(gomock.Any()).Return(nil, errors.New("not found"))
			},
			expectedErr: "track not found",
		},
		{
			name:    "track too short",
			userID:  userID,
			trackID: trackID,
			mockSetup: func() {
				trackRepo.EXPECT().FindByID(trackID).Return(&models.Track{
					ID:       trackID,
					Duration: 20, // 20 секунд
				}, nil)
			},
			expectedErr: "track is too short to record playback",
		},
		{
			name:    "played too frequently",
			userID:  userID,
			trackID: trackID,
			mockSetup: func() {
				trackRepo.EXPECT().FindByID(trackID).Return(&models.Track{
					ID:       trackID,
					Duration: 180,
				}, nil)
				historyRepo.EXPECT().GetHistory(userID).Return([]*models.ListeningHistory{
					{TrackID: trackID, ListenedAt: now.Add(-4 * time.Minute)},
					{TrackID: trackID, ListenedAt: now.Add(-3 * time.Minute)},
					{TrackID: trackID, ListenedAt: now.Add(-2 * time.Minute)},
				}, nil)
			},
			expectedErr: "track played too frequently",
		},
		{
			name:    "failed to record",
			userID:  userID,
			trackID: trackID,
			mockSetup: func() {
				trackRepo.EXPECT().FindByID(trackID).Return(&models.Track{
					ID:       trackID,
					Duration: 180,
				}, nil)
				historyRepo.EXPECT().GetHistory(userID).Return([]*models.ListeningHistory{}, nil)
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

			err := uc.RecordPlayback(tt.userID, tt.trackID)

			if tt.expectedErr != "" {
				assert.ErrorContains(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetUserHistory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	historyRepo := mocks.NewMockHistoryRepository(ctrl)
	trackRepo := mocks.NewMockTrackRepository(ctrl)
	uc := usecases.NewHistoryUseCase(historyRepo, trackRepo)

	userID := uuid.MustParse("1f22acbc-d723-4f92-a06f-8a583b4d45ea")
	now := time.Now()

	expectedTrackIDs := []uuid.UUID{
		uuid.MustParse("1e872412-8a31-4790-b118-ce335a3f4e84"),
		uuid.MustParse("2de4e9b3-dd50-457a-9f23-dcf83930cc17"),
		uuid.MustParse("168335ad-a5f0-44f7-ab1d-bb74398c6e37"),
	}

	tests := []struct {
		name        string
		userID      uuid.UUID
		mockSetup   func()
		expected    []*models.ListeningHistory
		expectedErr string
	}{
		{
			name:   "limit to 100 items",
			userID: userID,
			mockSetup: func() {
				history := make([]*models.ListeningHistory, 150)
				for i := range history {
					history[i] = &models.ListeningHistory{
						TrackID:    expectedTrackIDs[i%len(expectedTrackIDs)],
						ListenedAt: now.Add(-time.Duration(i) * time.Minute),
					}
				}
				historyRepo.EXPECT().GetHistory(userID).Return(history, nil)
			},
			expected: func() []*models.ListeningHistory {
				expected := make([]*models.ListeningHistory, 100)
				for i := 0; i < 100; i++ {
					expected[i] = &models.ListeningHistory{
						TrackID:    expectedTrackIDs[i%len(expectedTrackIDs)],
						ListenedAt: now.Add(-time.Duration(i) * time.Minute),
					}
				}
				return expected
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			result, err := uc.GetUserHistory(tt.userID)

			if tt.expectedErr != "" {
				assert.ErrorContains(t, err, tt.expectedErr)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, len(tt.expected))

				for i := 0; i < len(result)-1; i++ {
					assert.True(t, result[i].ListenedAt.After(result[i+1].ListenedAt),
						"history is not sorted correctly (newest first)")
				}
				assert.LessOrEqual(t, len(result), 100)
				for i := range tt.expected {
					assert.Equal(t, tt.expected[i].TrackID, result[i].TrackID)
					assert.Equal(t, tt.expected[i].ListenedAt.Unix(), result[i].ListenedAt.Unix())
				}
			}
		})
	}
}

func TestGetRecentPlays(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	historyRepo := mocks.NewMockHistoryRepository(ctrl)
	trackRepo := mocks.NewMockTrackRepository(ctrl)
	uc := usecases.NewHistoryUseCase(historyRepo, trackRepo)

	userID := uuid.MustParse("1f22acbc-d723-4f92-a06f-8a583b4d45ea")
	now := time.Now()

	tests := []struct {
		name      string
		userID    uuid.UUID
		within    time.Duration
		mockSetup func()
		expected  []*models.ListeningHistory
	}{
		{
			name:   "recent plays within 1 hour",
			userID: userID,
			within: time.Hour,
			mockSetup: func() {
				historyRepo.EXPECT().GetHistory(userID).Return([]*models.ListeningHistory{
					{TrackID: uuid.New(), ListenedAt: now.Add(-2 * time.Hour)},
					{TrackID: uuid.New(), ListenedAt: now.Add(-30 * time.Minute)},
					{TrackID: uuid.New(), ListenedAt: now.Add(-10 * time.Minute)},
				}, nil)
			},
			expected: []*models.ListeningHistory{
				{TrackID: uuid.New(), ListenedAt: now.Add(-30 * time.Minute)},
				{TrackID: uuid.New(), ListenedAt: now.Add(-10 * time.Minute)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			result, err := uc.GetRecentPlays(tt.userID, tt.within)

			assert.NoError(t, err)
			assert.Equal(t, len(tt.expected), len(result))

			for _, entry := range result {
				assert.True(t, now.Sub(entry.ListenedAt) <= tt.within,
					"entry is not within specified time range")
			}
		})
	}
}
