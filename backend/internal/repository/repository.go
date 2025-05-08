package repository

import (
	"database/sql"
	"music-service/internal/repository/db"
	"music-service/internal/repository/interfaces"
	"music-service/internal/repository/postgres"
)

type Repository struct {
	User     interfaces.UserRepository
	Track    interfaces.TrackRepository
	Album    interfaces.AlbumRepository
	Playlist interfaces.PlaylistRepository
	Genre    interfaces.GenreRepository
	Session  interfaces.SessionRepository
	History  interfaces.HistoryRepository
}

func NewRepository(cfg db.Config) (*Repository, error) {
	db, err := db.NewPostgresDB(cfg)
	if err != nil {
		return nil, err
	}

	return &Repository{
		User:     postgres.NewUserRepository(db),
		Track:    postgres.NewTrackRepository(db, cfg.TracksDir),
		Album:    postgres.NewAlbumRepository(db),
		Playlist: postgres.NewPlaylistRepository(db),
		Genre:    postgres.NewGenreRepository(db),
		Session:  postgres.NewSessionRepository(db),
		History:  postgres.NewHistoryRepository(db),
	}, nil
}

func NewRepositoryWithDB(db *sql.DB, tracksDir string) *Repository {
	return &Repository{
		User:     postgres.NewUserRepository(db),
		Track:    postgres.NewTrackRepository(db, tracksDir),
		Album:    postgres.NewAlbumRepository(db),
		Playlist: postgres.NewPlaylistRepository(db),
		Genre:    postgres.NewGenreRepository(db),
		Session:  postgres.NewSessionRepository(db),
		History:  postgres.NewHistoryRepository(db),
	}
}
