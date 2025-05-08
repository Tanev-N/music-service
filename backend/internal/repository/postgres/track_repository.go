package postgres

import (
	"database/sql"
	"music-service/internal/models"
	"music-service/internal/repository/interfaces"

	"github.com/google/uuid"
)

type TrackRepository struct {
	db *sql.DB
}

func NewTrackRepository(db *sql.DB) interfaces.TrackRepository {
	return &TrackRepository{
		db: db,
	}
}

func (r *TrackRepository) FindByID(id uuid.UUID) (*models.Track, error) {
	var track models.Track
	query := `SELECT id, title, duration, file_path, album_id, artist_name, cover_url, added_date, updated_at, play_count 
				FROM tracks WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(
		&track.ID,
		&track.Title,
		&track.Duration,
		&track.FilePath,
		&track.AlbumID,
		&track.ArtistName,
		&track.CoverURL,
		&track.AddedDate,
		&track.UpdatedAt,
		&track.PlayCount,
	)
	if err != nil {
		return nil, err
	}
	return &track, nil
}

func (r *TrackRepository) Save(track *models.Track) error {
	query := `
		INSERT INTO tracks (id, title, duration, file_path, album_id, artist_name, cover_url, added_date, updated_at, play_count) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (id) DO UPDATE 
		SET title = $2, duration = $3, file_path = $4, album_id = $5, artist_name = $6, 
			cover_url = $7, updated_at = $9, play_count = $10
	`
	_, err := r.db.Exec(query, track.ID, track.Title, track.Duration, track.FilePath,
		track.AlbumID, track.ArtistName, track.CoverURL, track.AddedDate, track.UpdatedAt, track.PlayCount)
	return err
}

func (r *TrackRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM tracks WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *TrackRepository) Search(query string) ([]*models.Track, error) {
	var tracks []*models.Track
	rows, err := r.db.Query(`SELECT id, title, duration, file_path, album_id, artist_name, cover_url, added_date, updated_at, play_count 
					FROM tracks WHERE title ILIKE $1 OR artist_name ILIKE $1`, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var track models.Track
		err := rows.Scan(
			&track.ID,
			&track.Title,
			&track.Duration,
			&track.FilePath,
			&track.AlbumID,
			&track.ArtistName,
			&track.CoverURL,
			&track.AddedDate,
			&track.UpdatedAt,
			&track.PlayCount,
		)
		if err != nil {
			return nil, err
		}
		tracks = append(tracks, &track)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tracks, nil
}

func (r *TrackRepository) IncrementPlayCount(trackID uuid.UUID) error {
	query := `UPDATE tracks SET play_count = play_count + 1 WHERE id = $1`
	_, err := r.db.Exec(query, trackID)
	return err
}
