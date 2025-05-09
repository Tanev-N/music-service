package postgres

import (
	"database/sql"
	"fmt"
	"io"
	"music-service/internal/models"
	"music-service/internal/repository/interfaces"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
)

type TrackRepository struct {
	db        *sql.DB
	tracksDir string
}

func NewTrackRepository(db *sql.DB, tracksDir string) interfaces.TrackRepository {
	return &TrackRepository{
		db:        db,
		tracksDir: tracksDir,
	}
}

func (r *TrackRepository) FindByID(id uuid.UUID) (*models.Track, error) {
	var track models.Track
	var albumID pgtype.UUID
	query := `SELECT id, title, duration, file_path, album_id, artist_name, cover_url, added_date, updated_at, play_count 
				FROM tracks WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(
		&track.ID,
		&track.Title,
		&track.Duration,
		&track.FilePath,
		&albumID,
		&track.ArtistName,
		&track.CoverURL,
		&track.AddedDate,
		&track.UpdatedAt,
		&track.PlayCount,
	)
	if err != nil {
		return nil, err
	}

	if albumID.Status == pgtype.Present {
		track.AlbumID = albumID.Bytes
	}

	return &track, nil
}

func (r *TrackRepository) Save(track *models.Track) error {
	var albumID interface{}
	if track.AlbumID == uuid.Nil {
		albumID = nil
	} else {
		albumID = track.AlbumID
	}

	query := `
		INSERT INTO tracks (id, title, duration, file_path, album_id, artist_name, cover_url, added_date, updated_at, play_count) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (id) DO UPDATE 
		SET title = $2, duration = $3, file_path = $4, album_id = $5, artist_name = $6, 
			cover_url = $7, updated_at = $9, play_count = $10
	`
	_, err := r.db.Exec(query, track.ID, track.Title, track.Duration, track.FilePath,
		albumID, track.ArtistName, track.CoverURL, track.AddedDate, track.UpdatedAt, track.PlayCount)
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

func (r *TrackRepository) SaveTrackFile(trackID uuid.UUID, fileReader io.Reader, fileSize int64) (string, error) {
	if err := os.MkdirAll(r.tracksDir, 0755); err != nil {
		return "", fmt.Errorf("не удалось создать директорию для треков: %w", err)
	}

	trackIDStr := trackID.String()
	relativeDir := filepath.Join(trackIDStr[:2], trackIDStr[2:4])
	absoluteDir := filepath.Join(r.tracksDir, relativeDir)

	if err := os.MkdirAll(absoluteDir, 0755); err != nil {
		return "", fmt.Errorf("не удалось создать поддиректорию: %w", err)
	}

	fileName := fmt.Sprintf("%s.mp3", trackIDStr)
	absolutePath := filepath.Join(absoluteDir, fileName)
	relativePath := filepath.Join(relativeDir, fileName)

	file, err := os.Create(absolutePath)
	if err != nil {
		return "", fmt.Errorf("не удалось создать файл: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, fileReader)
	if err != nil {
		os.Remove(absolutePath)
		return "", fmt.Errorf("не удалось сохранить файл: %w", err)
	}

	return relativePath, nil
}

func (r *TrackRepository) GetStorageDir() string {
	return r.tracksDir
}

// GetGenresForTrack возвращает список жанров для трека
func (r *TrackRepository) GetGenresForTrack(trackID uuid.UUID) ([]*models.Genre, error) {
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
