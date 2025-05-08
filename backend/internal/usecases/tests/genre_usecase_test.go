package usecases_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"music-service/internal/models"
	"music-service/internal/repository/interfaces/mocks"
	"music-service/internal/usecases"
)

func TestCreateGenre(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	genreRepo := mocks.NewMockGenreRepository(ctrl)
	trackRepo := mocks.NewMockTrackRepository(ctrl)
	uc := usecases.NewGenreUseCase(genreRepo, trackRepo)

	// Фиксированный UUID для жанра
	genreID := uuid.MustParse("1f22acbc-d723-4f92-a06f-8a583b4d45ea")

	tests := []struct {
		name        string
		genreName   string
		mockSetup   func()
		expected    *models.Genre
		expectedErr string
	}{
		{
			name:      "successful creation",
			genreName: "Rock",
			mockSetup: func() {
				genreRepo.EXPECT().ListAll().Return([]*models.Genre{}, nil)
				genreRepo.EXPECT().Save(gomock.Any()).DoAndReturn(func(g *models.Genre) error {
					g.ID = genreID // Используем фиксированный UUID
					return nil
				})
			},
			expected: &models.Genre{
				ID:   genreID, // Используем фиксированный UUID
				Name: "Rock",
			},
		},
		{
			name:        "name too short",
			genreName:   "A",
			expectedErr: "genre name must be at least 2 characters",
		},
		{
			name:        "name too long",
			genreName:   strings.Repeat("a", 51),
			expectedErr: "genre name is too long (max 50 characters)",
		},
		{
			name:      "genre already exists",
			genreName: "Existing Genre",
			mockSetup: func() {
				genreRepo.EXPECT().ListAll().Return([]*models.Genre{
					{Name: "Existing Genre"},
				}, nil)
			},
			expectedErr: "genre already exists",
		},
		{
			name:      "failed to save",
			genreName: "Rock",
			mockSetup: func() {
				genreRepo.EXPECT().ListAll().Return([]*models.Genre{}, nil)
				genreRepo.EXPECT().Save(gomock.Any()).Return(errors.New("save failed"))
			},
			expectedErr: "failed to save genre",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			result, err := uc.CreateGenre(tt.genreName)

			if tt.expectedErr != "" {
				assert.ErrorContains(t, err, tt.expectedErr)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.Name, result.Name)
				assert.Equal(t, tt.expected.ID, result.ID) // Проверяем фиксированный UUID
			}
		})
	}
}

func TestGetGenresByTrack(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	genreRepo := mocks.NewMockGenreRepository(ctrl)
	trackRepo := mocks.NewMockTrackRepository(ctrl)
	uc := usecases.NewGenreUseCase(genreRepo, trackRepo)

	trackID := uuid.New()

	tests := []struct {
		name        string
		trackID     uuid.UUID
		mockSetup   func()
		expected    []*models.Genre
		expectedErr string
	}{
		{
			name:    "successful get",
			trackID: trackID,
			mockSetup: func() {
				trackRepo.EXPECT().FindByID(trackID).Return(&models.Track{ID: trackID}, nil)
				genreRepo.EXPECT().GetGenresForTrack(trackID).Return([]*models.Genre{
					{ID: uuid.MustParse("1f22acbc-d723-4f92-a06f-8a583b4d45ea"), Name: "Rock"},
				}, nil)
			},
			expected: []*models.Genre{
				{ID: uuid.MustParse("1f22acbc-d723-4f92-a06f-8a583b4d45ea"), Name: "Rock"},
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
			name:    "failed to get genres",
			trackID: trackID,
			mockSetup: func() {
				trackRepo.EXPECT().FindByID(trackID).Return(&models.Track{ID: trackID}, nil)
				genreRepo.EXPECT().GetGenresForTrack(trackID).Return(nil, errors.New("failed"))
			},
			expectedErr: "failed to get track genres",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			result, err := uc.GetGenresByTrack(tt.trackID)

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

func TestListAllGenres(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	genreRepo := mocks.NewMockGenreRepository(ctrl)
	trackRepo := mocks.NewMockTrackRepository(ctrl)
	uc := usecases.NewGenreUseCase(genreRepo, trackRepo)

	rockID := uuid.MustParse("1f22acbc-d723-4f92-a06f-8a583b4d45ea")
	popID := uuid.MustParse("2f46e5ce-8b10-4994-bc10-d33637e98ebf")

	tests := []struct {
		name        string
		mockSetup   func()
		expected    []*models.Genre
		expectedErr string
	}{
		{
			name: "successful list",
			mockSetup: func() {
				genreRepo.EXPECT().ListAll().Return([]*models.Genre{
					{ID: rockID, Name: "Rock"},
					{ID: popID, Name: "Pop"},
				}, nil)
			},
			expected: []*models.Genre{
				{ID: popID, Name: "Pop"},
				{ID: rockID, Name: "Rock"},
			},
		},
		{
			name: "failed to list",
			mockSetup: func() {
				genreRepo.EXPECT().ListAll().Return(nil, errors.New("failed"))
			},
			expectedErr: "failed to list genres",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			result, err := uc.ListAllGenres()

			if tt.expectedErr != "" {
				assert.ErrorContains(t, err, tt.expectedErr)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
				for i := 0; i < len(result)-1; i++ {
					assert.True(t, strings.Compare(result[i].Name, result[i+1].Name) <= 0)
				}
			}
		})
	}
}

func TestRemoveGenreFromTrack(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	genreRepo := mocks.NewMockGenreRepository(ctrl)
	trackRepo := mocks.NewMockTrackRepository(ctrl)
	uc := usecases.NewGenreUseCase(genreRepo, trackRepo)

	trackID := uuid.New()
	genreID := uuid.MustParse("4e5f6a7b-8c9d-4e0f-a1b2-3c4d5e6f7a8b") // Фиксированный UUID для жанра

	tests := []struct {
		name        string
		trackID     uuid.UUID
		genreID     uuid.UUID
		mockSetup   func()
		expectedErr string
	}{
		{
			name:    "successful remove",
			trackID: trackID,
			genreID: genreID,
			mockSetup: func() {
				genreRepo.EXPECT().GetGenresForTrack(trackID).Return([]*models.Genre{
					{ID: genreID},
				}, nil)
				genreRepo.EXPECT().RemoveGenreFromTrack(trackID, genreID).Return(nil)
			},
		},
		{
			name:    "genre not assigned",
			trackID: trackID,
			genreID: uuid.New(),
			mockSetup: func() {
				genreRepo.EXPECT().GetGenresForTrack(trackID).Return([]*models.Genre{
					{ID: genreID},
				}, nil)
			},
			expectedErr: "genre is not assigned to this track",
		},
		{
			name:    "failed to remove",
			trackID: trackID,
			genreID: genreID,
			mockSetup: func() {
				genreRepo.EXPECT().GetGenresForTrack(trackID).Return([]*models.Genre{
					{ID: genreID},
				}, nil)
				genreRepo.EXPECT().RemoveGenreFromTrack(trackID, genreID).Return(errors.New("failed"))
			},
			expectedErr: "failed to remove genre from track",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			err := uc.RemoveGenreFromTrack(tt.trackID, tt.genreID)

			if tt.expectedErr != "" {
				assert.ErrorContains(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
