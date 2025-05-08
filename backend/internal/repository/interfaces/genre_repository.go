package interfaces

import (
	"music-service/internal/models"

	"github.com/google/uuid"
)

type GenreRepository interface {
	FindByID(id uuid.UUID) (*models.Genre, error)
	Save(genre *models.Genre) error
	Delete(id uuid.UUID) error
	GetGenresForTrack(trackID uuid.UUID) ([]*models.Genre, error)
	AddGenreToTrack(trackID, genreID uuid.UUID) error
	RemoveGenreFromTrack(trackID, genreID uuid.UUID) error
	ListAll() ([]*models.Genre, error)
}
