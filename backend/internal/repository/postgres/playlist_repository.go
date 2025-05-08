package postgres

import (
	"database/sql"
	"music-service/internal/models"
	"music-service/internal/repository/interfaces"

	"github.com/google/uuid"
)

type PlaylistRepository struct {
	db *sql.DB
}

func NewPlaylistRepository(db *sql.DB) interfaces.PlaylistRepository {
	return &PlaylistRepository{
		db: db,
	}
}

func (r *PlaylistRepository) FindByID(id uuid.UUID) (*models.Playlist, error) {
	var playlist models.Playlist
	query := `SELECT id, name, description, user_id, cover_url, created_date, updated_at FROM playlists WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(
		&playlist.ID,
		&playlist.Name,
		&playlist.Description,
		&playlist.UserID,
		&playlist.CoverURL,
		&playlist.CreatedDate,
		&playlist.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &playlist, nil
}

func (r *PlaylistRepository) Save(playlist *models.Playlist) error {
	query := `
		INSERT INTO playlists (id, name, description, user_id, cover_url, created_date, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE 
		SET name = $2, description = $3, cover_url = $5, updated_at = $7
	`
	_, err := r.db.Exec(query,
		playlist.ID,
		playlist.Name,
		playlist.Description,
		playlist.UserID,
		playlist.CoverURL,
		playlist.CreatedDate,
		playlist.UpdatedAt,
	)
	return err
}

func (r *PlaylistRepository) Delete(id uuid.UUID) error {
	// Сначала удаляем связи с треками
	_, err := r.db.Exec(`DELETE FROM playlist_tracks WHERE playlist_id = $1`, id)
	if err != nil {
		return err
	}
	// Затем удаляем сам плейлист
	_, err = r.db.Exec(`DELETE FROM playlists WHERE id = $1`, id)
	return err
}

func (r *PlaylistRepository) AddTrack(playlistID uuid.UUID, trackID uuid.UUID) error {
	query := `
		INSERT INTO playlist_tracks (playlist_id, track_id, added_at) 
		VALUES ($1, $2, NOW())
		ON CONFLICT (playlist_id, track_id) DO NOTHING
	`
	_, err := r.db.Exec(query, playlistID, trackID)
	return err
}

func (r *PlaylistRepository) RemoveTrack(playlistID uuid.UUID, trackID uuid.UUID) error {
	query := `DELETE FROM playlist_tracks WHERE playlist_id = $1 AND track_id = $2`
	_, err := r.db.Exec(query, playlistID, trackID)
	return err
}

func (r *PlaylistRepository) GetTracks(playlistID uuid.UUID) ([]*models.Track, error) {
	var tracks []*models.Track
	query := `
		SELECT t.id, t.title, t.duration, t.file_path, t.album_id, t.artist_name, t.cover_url, t.added_date, t.updated_at, t.play_count
		FROM tracks t
		JOIN playlist_tracks pt ON t.id = pt.track_id
		WHERE pt.playlist_id = $1
		ORDER BY pt.added_at DESC
	`
	rows, err := r.db.Query(query, playlistID)
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

func (r *PlaylistRepository) GetUserPlaylists(userID uuid.UUID) ([]*models.Playlist, error) {
	var playlists []*models.Playlist
	query := `SELECT id, name, description, user_id, cover_url, created_date, updated_at FROM playlists WHERE user_id = $1`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var playlist models.Playlist
		err := rows.Scan(
			&playlist.ID,
			&playlist.Name,
			&playlist.Description,
			&playlist.UserID,
			&playlist.CoverURL,
			&playlist.CreatedDate,
			&playlist.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		playlists = append(playlists, &playlist)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return playlists, nil
}
