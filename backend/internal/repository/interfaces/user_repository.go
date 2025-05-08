package interfaces

import (
	"music-service/internal/models"

	"github.com/google/uuid"
)

type UserRepository interface {
	FindByID(id uuid.UUID) (*models.User, error)
	Save(user *models.User) error
	Delete(id uuid.UUID) error
	Search(query string) ([]*models.User, error)
}
