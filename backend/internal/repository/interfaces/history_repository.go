package interfaces

import (
	"music-service/internal/models"

	"github.com/google/uuid"
)

type HistoryRepository interface {
	AddEntry(userID uuid.UUID, trackID uuid.UUID) error
	GetHistory(userID uuid.UUID) ([]*models.ListeningHistory, error)
	GetPlayCount(trackID uuid.UUID) (int, error)
}
