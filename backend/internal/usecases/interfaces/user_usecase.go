package interfaces

import (
	"music-service/internal/models"

	"github.com/google/uuid"
)

type UserUseCase interface {
	Register(login, password string) (*models.User, error)
	Authenticate(login, password string) (*models.User, *models.Session, error)
	GetUserProfile(userID uuid.UUID) (*models.User, error)
	UpdatePermissions(userID uuid.UUID, permission models.Permission) error
	DeleteUser(userID uuid.UUID) error
	Logout(sessionID uuid.UUID) error
	ValidateSession(token string) (*models.User, error)
}
