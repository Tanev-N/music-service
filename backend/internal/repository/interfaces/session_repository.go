package interfaces

import (
	"music-service/internal/models"

	"github.com/google/uuid"
)

type SessionRepository interface {
	CreateSession(userID uuid.UUID, token string) (*models.Session, error)
	GetSession(sessionID string) (*models.Session, error)
	DeleteSession(sessionID string) error
	DeleteAllForUser(userID uuid.UUID) error
	GetSessionByToken(token string) (*models.Session, *models.User, error)
}
