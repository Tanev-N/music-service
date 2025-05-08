package interfaces

import (
	"music-service/internal/models"
	"time"

	"github.com/google/uuid"
)

type HistoryUseCase interface {
	RecordPlayback(userID uuid.UUID, trackID uuid.UUID) error
	GetUserHistory(userID uuid.UUID) ([]*models.ListeningHistory, error)
	GetRecentPlays(userID uuid.UUID, within time.Duration) ([]*models.ListeningHistory, error)
}
