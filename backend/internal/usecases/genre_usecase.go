package usecases

import (
	"errors"
	"fmt"
	"music-service/internal/models"
	"music-service/internal/repository/interfaces"
	usecaseInterfaces "music-service/internal/usecases/interfaces"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
)

type genreUseCase struct {
	genreRepo interfaces.GenreRepository
	trackRepo interfaces.TrackRepository
}

func NewGenreUseCase(
	genreRepo interfaces.GenreRepository,
	trackRepo interfaces.TrackRepository,
) usecaseInterfaces.GenreUseCase {
	return &genreUseCase{
		genreRepo: genreRepo,
		trackRepo: trackRepo,
	}
}

func (uc *genreUseCase) CreateGenre(name string) (*models.Genre, error) {
	name = strings.TrimSpace(name)
	if utf8.RuneCountInString(name) < 2 {
		return nil, errors.New("genre name must be at least 2 characters")
	}
	if utf8.RuneCountInString(name) > 50 {
		return nil, errors.New("genre name is too long (max 50 characters)")
	}
	existingGenres, err := uc.genreRepo.ListAll()
	if err != nil {
		return nil, fmt.Errorf("failed to check existing genres: %w", err)
	}
	for _, g := range existingGenres {
		if strings.EqualFold(g.Name, name) {
			return nil, errors.New("genre already exists")
		}
	}
	genre := &models.Genre{
		ID:   uuid.New(),
		Name: name,
	}
	if err := uc.genreRepo.Save(genre); err != nil {
		return nil, fmt.Errorf("failed to save genre: %w", err)
	}
	return genre, nil
}

func (uc *genreUseCase) GetGenresByTrack(trackID uuid.UUID) ([]*models.Genre, error) {
	if _, err := uc.trackRepo.FindByID(trackID); err != nil {
		return nil, fmt.Errorf("track not found: %w", err)
	}
	genres, err := uc.genreRepo.GetGenresForTrack(trackID)
	if err != nil {
		return nil, fmt.Errorf("failed to get track genres: %w", err)
	}
	return genres, nil
}

func (uc *genreUseCase) ListAllGenres() ([]*models.Genre, error) {
	genres, err := uc.genreRepo.ListAll()
	if err != nil {
		return nil, fmt.Errorf("failed to list genres: %w", err)
	}
	// Сортируем жанры по алфавиту
	for i := 0; i < len(genres)-1; i++ {
		for j := i + 1; j < len(genres); j++ {
			if strings.Compare(genres[i].Name, genres[j].Name) > 0 {
				genres[i], genres[j] = genres[j], genres[i]
			}
		}
	}

	return genres, nil
}

func (uc *genreUseCase) AssignGenreToTrack(trackID, genreID uuid.UUID) error {
	if _, err := uc.trackRepo.FindByID(trackID); err != nil {
		return fmt.Errorf("track not found: %w", err)
	}
	if _, err := uc.genreRepo.FindByID(genreID); err != nil {
		return fmt.Errorf("genre not found: %w", err)
	}
	currentGenres, err := uc.genreRepo.GetGenresForTrack(trackID)
	if err != nil {
		return fmt.Errorf("failed to get track genres: %w", err)
	}
	for _, g := range currentGenres {
		if g.ID == genreID {
			return errors.New("genre already assigned to track")
		}
	}
	if len(currentGenres) >= 5 {
		return errors.New("track cannot have more than 5 genres")
	}

	if err := uc.genreRepo.AddGenreToTrack(trackID, genreID); err != nil {
		return fmt.Errorf("failed to assign genre to track: %w", err)
	}

	return nil
}

func (uc *genreUseCase) RemoveGenreFromTrack(trackID, genreID uuid.UUID) error {
	genres, err := uc.genreRepo.GetGenresForTrack(trackID)
	if err != nil {
		return fmt.Errorf("failed to get track genres: %w", err)
	}

	found := false
	for _, g := range genres {
		if g.ID == genreID {
			found = true
			break
		}
	}

	if !found {
		return errors.New("genre is not assigned to this track")
	}

	if err := uc.genreRepo.RemoveGenreFromTrack(trackID, genreID); err != nil {
		return fmt.Errorf("failed to remove genre from track: %w", err)
	}

	return nil
}
