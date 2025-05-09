package postgres

import (
	"database/sql"
	"music-service/internal/models"
	"music-service/internal/repository/interfaces"

	"github.com/google/uuid"
)

type GenreRepository struct {
	db *sql.DB
}

func NewGenreRepository(db *sql.DB) interfaces.GenreRepository {
	return &GenreRepository{
		db: db,
	}
}

func (r *GenreRepository) FindByID(id uuid.UUID) (*models.Genre, error) {
	var genre models.Genre
	query := `SELECT id, name FROM genres WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&genre.ID, &genre.Name)
	if err != nil {
		return nil, err
	}
	return &genre, nil
}

func (r *GenreRepository) Save(genre *models.Genre) error {
	query := `
		INSERT INTO genres (id, name) 
		VALUES ($1, $2)
		ON CONFLICT (id) DO UPDATE 
		SET name = $2
	`
	_, err := r.db.Exec(query, genre.ID, genre.Name)
	return err
}

func (r *GenreRepository) Delete(id uuid.UUID) error {
	_, err := r.db.Exec(`DELETE FROM track_genres WHERE genre_id = $1`, id)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(`DELETE FROM genres WHERE id = $1`, id)
	return err
}

func (r *GenreRepository) GetGenresForTrack(trackID uuid.UUID) ([]*models.Genre, error) {
	var genres []*models.Genre
	query := `
		SELECT g.id, g.name
		FROM genres g
		JOIN track_genres tg ON g.id = tg.genre_id
		WHERE tg.track_id = $1
	`
	rows, err := r.db.Query(query, trackID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var genre models.Genre
		err := rows.Scan(&genre.ID, &genre.Name)
		if err != nil {
			return nil, err
		}
		genres = append(genres, &genre)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return genres, nil
}

func (r *GenreRepository) AddGenreToTrack(trackID, genreID uuid.UUID) error {
	query := `
		INSERT INTO track_genres (track_id, genre_id) 
		VALUES ($1, $2)
		ON CONFLICT (track_id, genre_id) DO NOTHING
	`
	_, err := r.db.Exec(query, trackID, genreID)
	return err
}

func (r *GenreRepository) RemoveGenreFromTrack(trackID, genreID uuid.UUID) error {
	query := `DELETE FROM track_genres WHERE track_id = $1 AND genre_id = $2`
	_, err := r.db.Exec(query, trackID, genreID)
	return err
}

func (r *GenreRepository) ListAll() ([]*models.Genre, error) {
	var genres []*models.Genre
	query := `SELECT id, name FROM genres`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var genre models.Genre
		err := rows.Scan(&genre.ID, &genre.Name)
		if err != nil {
			return nil, err
		}
		genres = append(genres, &genre)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return genres, nil
}
