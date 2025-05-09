package models

import (
	"time"

	"github.com/google/uuid"
)

type ListeningHistory struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	TrackID    uuid.UUID
	ListenedAt time.Time
	Track      Track
}
