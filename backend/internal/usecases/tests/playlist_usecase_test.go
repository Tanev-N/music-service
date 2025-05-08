package usecases_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"music-service/internal/models"
	"music-service/internal/repository/interfaces/mocks"
	"music-service/internal/usecases"
)

func TestRemoveTrackFromPlaylist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	playlistRepo := mocks.NewMockPlaylistRepository(ctrl)
	trackRepo := mocks.NewMockTrackRepository(ctrl)
	userRepo := mocks.NewMockUserRepository(ctrl)
	uc := usecases.NewPlaylistUseCase(playlistRepo, trackRepo, userRepo)

	playlistID := uuid.New()
	trackID := uuid.New()

	t.Run("successful track removal", func(t *testing.T) {
		playlistRepo.EXPECT().FindByID(playlistID).Return(&models.Playlist{ID: playlistID}, nil)
		playlistRepo.EXPECT().GetTracks(playlistID).Return([]*models.Track{{ID: trackID}}, nil)
		playlistRepo.EXPECT().RemoveTrack(playlistID, trackID).Return(nil)
		playlistRepo.EXPECT().Save(gomock.Any()).Return(nil)

		err := uc.RemoveTrackFromPlaylist(playlistID, trackID)
		assert.NoError(t, err)
	})
}
