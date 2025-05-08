package interfaces

import (
	"music-service/internal/models"

	"github.com/google/uuid"
)

type GenreUseCase interface {
	CreateGenre(name string) (*models.Genre, error)
	GetGenresByTrack(trackID uuid.UUID) ([]*models.Genre, error)
	ListAllGenres() ([]*models.Genre, error)
	RemoveGenreFromTrack(trackID uuid.UUID, genreID uuid.UUID) error
	AssignGenreToTrack(trackID uuid.UUID, genreID uuid.UUID) error
}
