package postgres

import (
	"database/sql"
	"music-service/internal/models"
	"music-service/internal/repository/interfaces"

	"github.com/google/uuid"
)

type AlbumRepository struct {
	db *sql.DB
}

func NewAlbumRepository(db *sql.DB) interfaces.AlbumRepository {
	return &AlbumRepository{
		db: db,
	}
}

func (r *AlbumRepository) FindByID(id uuid.UUID) (*models.Album, error) {
	var album models.Album
	query := `SELECT id, title, release_date, cover_url, created_at, updated_at FROM albums WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(
		&album.ID,
		&album.Title,
		&album.ReleaseDate,
		&album.CoverURL,
		&album.CreatedAt,
		&album.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &album, nil
}

func (r *AlbumRepository) Save(album *models.Album) error {
	query := `
		INSERT INTO albums (id, title, release_date, cover_url, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE 
		SET title = $2, release_date = $3, cover_url = $4, updated_at = $6
	`
	_, err := r.db.Exec(query,
		album.ID,
		album.Title,
		album.ReleaseDate,
		album.CoverURL,
		album.CreatedAt,
		album.UpdatedAt,
	)
	return err
}

func (r *AlbumRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM albums WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *AlbumRepository) GetTracks(albumID uuid.UUID) ([]*models.Track, error) {
	var tracks []*models.Track
	query := `SELECT id, title, duration, file_path, album_id, artist_name, cover_url, added_date, updated_at, play_count 
				FROM tracks WHERE album_id = $1`

	rows, err := r.db.Query(query, albumID)
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

func (r *AlbumRepository) AddTrackToAlbum(albumID, trackID uuid.UUID) error {
	query := `UPDATE tracks SET album_id = $1 WHERE id = $2`
	_, err := r.db.Exec(query, albumID, trackID)
	return err
}

func (r *AlbumRepository) RemoveTrackFromAlbum(albumID, trackID uuid.UUID) error {
	query := `UPDATE tracks SET album_id = NULL WHERE id = $2 AND album_id = $1`
	_, err := r.db.Exec(query, albumID, trackID)
	return err
}

func (r *AlbumRepository) ListAll() ([]*models.Album, error) {
	var albums []*models.Album
	query := `SELECT id, title, release_date, cover_url, created_at, updated_at FROM albums`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var album models.Album
		err := rows.Scan(
			&album.ID,
			&album.Title,
			&album.ReleaseDate,
			&album.CoverURL,
			&album.CreatedAt,
			&album.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		albums = append(albums, &album)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return albums, nil
}
