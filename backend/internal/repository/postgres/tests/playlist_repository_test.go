package tests

import (
	"errors"
	"music-service/internal/models"
	"music-service/internal/repository/postgres"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPlaylistRepository_FindByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewPlaylistRepository(db)

	playlistID := uuid.New()
	userID := uuid.New()
	now := time.Now()
	playlist := &models.Playlist{
		ID:          playlistID,
		Name:        "Test Playlist",
		Description: "Test Description",
		UserID:      userID,
		CoverURL:    "http://example.com/cover.jpg",
		CreatedDate: now,
		UpdatedAt:   now,
	}

	// Успешный сценарий
	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "description", "user_id", "cover_url", "created_date", "updated_at"}).
			AddRow(playlist.ID, playlist.Name, playlist.Description, playlist.UserID, playlist.CoverURL, playlist.CreatedDate, playlist.UpdatedAt)

		mock.ExpectQuery("SELECT (.+) FROM playlists WHERE id = ?").
			WithArgs(playlistID).
			WillReturnRows(rows)

		foundPlaylist, err := repo.FindByID(playlistID)
		assert.NoError(t, err)
		assert.Equal(t, playlist.ID, foundPlaylist.ID)
		assert.Equal(t, playlist.Name, foundPlaylist.Name)
		assert.Equal(t, playlist.Description, foundPlaylist.Description)
		assert.Equal(t, playlist.UserID, foundPlaylist.UserID)
	})

	// Сценарий с ошибкой
	t.Run("error", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM playlists WHERE id = ?").
			WithArgs(playlistID).
			WillReturnError(errors.New("db error"))

		foundPlaylist, err := repo.FindByID(playlistID)
		assert.Error(t, err)
		assert.Nil(t, foundPlaylist)
	})

	// Проверка, что все ожидаемые запросы были выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPlaylistRepository_Save(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewPlaylistRepository(db)

	playlistID := uuid.New()
	userID := uuid.New()
	now := time.Now()
	playlist := &models.Playlist{
		ID:          playlistID,
		Name:        "Test Playlist",
		Description: "Test Description",
		UserID:      userID,
		CoverURL:    "http://example.com/cover.jpg",
		CreatedDate: now,
		UpdatedAt:   now,
	}

	// Успешное сохранение
	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO playlists").
			WithArgs(playlist.ID, playlist.Name, playlist.Description, playlist.UserID, playlist.CoverURL, playlist.CreatedDate, playlist.UpdatedAt).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Save(playlist)
		assert.NoError(t, err)
	})

	// Ошибка при сохранении
	t.Run("error", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO playlists").
			WithArgs(playlist.ID, playlist.Name, playlist.Description, playlist.UserID, playlist.CoverURL, playlist.CreatedDate, playlist.UpdatedAt).
			WillReturnError(errors.New("db error"))

		err := repo.Save(playlist)
		assert.Error(t, err)
	})

	// Проверка, что все ожидаемые запросы были выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPlaylistRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewPlaylistRepository(db)

	playlistID := uuid.New()

	// Успешное удаление
	t.Run("success", func(t *testing.T) {
		// Сначала удаляются связи с треками
		mock.ExpectExec("DELETE FROM playlist_tracks WHERE playlist_id = ?").
			WithArgs(playlistID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Затем удаляется сам плейлист
		mock.ExpectExec("DELETE FROM playlists WHERE id = ?").
			WithArgs(playlistID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Delete(playlistID)
		assert.NoError(t, err)
	})

	// Ошибка при удалении связей
	t.Run("error_tracks", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM playlist_tracks WHERE playlist_id = ?").
			WithArgs(playlistID).
			WillReturnError(errors.New("db error"))

		err := repo.Delete(playlistID)
		assert.Error(t, err)
	})

	// Ошибка при удалении плейлиста
	t.Run("error_playlist", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM playlist_tracks WHERE playlist_id = ?").
			WithArgs(playlistID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectExec("DELETE FROM playlists WHERE id = ?").
			WithArgs(playlistID).
			WillReturnError(errors.New("db error"))

		err := repo.Delete(playlistID)
		assert.Error(t, err)
	})

	// Проверка, что все ожидаемые запросы были выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPlaylistRepository_GetTracks(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewPlaylistRepository(db)

	playlistID := uuid.New()
	now := time.Now()
	tracks := []*models.Track{
		{
			ID:         uuid.New(),
			Title:      "Track 1",
			Duration:   180,
			ArtistName: "Test Artist",
			AddedDate:  now,
			UpdatedAt:  now,
		},
		{
			ID:         uuid.New(),
			Title:      "Track 2",
			Duration:   240,
			ArtistName: "Test Artist",
			AddedDate:  now,
			UpdatedAt:  now,
		},
	}

	// Успешное получение треков
	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title", "duration", "file_path", "album_id", "artist_name", "cover_url", "added_date", "updated_at", "play_count"})
		for _, track := range tracks {
			rows.AddRow(track.ID, track.Title, track.Duration, track.FilePath, track.AlbumID, track.ArtistName, track.CoverURL, track.AddedDate, track.UpdatedAt, track.PlayCount)
		}

		mock.ExpectQuery("SELECT t.id, t.title, t.duration, t.file_path, t.album_id, t.artist_name, t.cover_url, t.added_date, t.updated_at, t.play_count FROM tracks t JOIN playlist_tracks pt ON t.id = pt.track_id WHERE pt.playlist_id = ?").
			WithArgs(playlistID).
			WillReturnRows(rows)

		foundTracks, err := repo.GetTracks(playlistID)
		assert.NoError(t, err)
		assert.Len(t, foundTracks, 2)
		assert.Equal(t, tracks[0].Title, foundTracks[0].Title)
		assert.Equal(t, tracks[1].Title, foundTracks[1].Title)
	})

	// Ошибка при получении треков
	t.Run("error", func(t *testing.T) {
		mock.ExpectQuery("SELECT t.id, t.title, t.duration, t.file_path, t.album_id, t.artist_name, t.cover_url, t.added_date, t.updated_at, t.play_count FROM tracks t JOIN playlist_tracks pt ON t.id = pt.track_id WHERE pt.playlist_id = ?").
			WithArgs(playlistID).
			WillReturnError(errors.New("db error"))

		foundTracks, err := repo.GetTracks(playlistID)
		assert.Error(t, err)
		assert.Nil(t, foundTracks)
	})

	// Проверка, что все ожидаемые запросы были выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPlaylistRepository_GetUserPlaylists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewPlaylistRepository(db)

	userID := uuid.New()
	now := time.Now()
	playlists := []*models.Playlist{
		{
			ID:          uuid.New(),
			Name:        "Playlist 1",
			Description: "Description 1",
			UserID:      userID,
			CreatedDate: now,
			UpdatedAt:   now,
		},
		{
			ID:          uuid.New(),
			Name:        "Playlist 2",
			Description: "Description 2",
			UserID:      userID,
			CreatedDate: now,
			UpdatedAt:   now,
		},
	}

	// Успешное получение плейлистов пользователя
	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "description", "user_id", "cover_url", "created_date", "updated_at"})
		for _, playlist := range playlists {
			rows.AddRow(playlist.ID, playlist.Name, playlist.Description, playlist.UserID, playlist.CoverURL, playlist.CreatedDate, playlist.UpdatedAt)
		}

		mock.ExpectQuery("SELECT id, name, description, user_id, cover_url, created_date, updated_at FROM playlists WHERE user_id = ?").
			WithArgs(userID).
			WillReturnRows(rows)

		foundPlaylists, err := repo.GetUserPlaylists(userID)
		assert.NoError(t, err)
		assert.Len(t, foundPlaylists, 2)
		assert.Equal(t, playlists[0].Name, foundPlaylists[0].Name)
		assert.Equal(t, playlists[1].Name, foundPlaylists[1].Name)
	})

	// Ошибка при получении плейлистов
	t.Run("error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, name, description, user_id, cover_url, created_date, updated_at FROM playlists WHERE user_id = ?").
			WithArgs(userID).
			WillReturnError(errors.New("db error"))

		foundPlaylists, err := repo.GetUserPlaylists(userID)
		assert.Error(t, err)
		assert.Nil(t, foundPlaylists)
	})

	// Проверка, что все ожидаемые запросы были выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
