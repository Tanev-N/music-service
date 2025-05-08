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

func TestAlbumRepository_FindByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewAlbumRepository(db)

	albumID := uuid.New()
	now := time.Now()
	album := &models.Album{
		ID:          albumID,
		Title:       "Test Album",
		ReleaseDate: now,
		CoverURL:    "http://example.com/cover.jpg",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Успешный сценарий
	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title", "release_date", "cover_url", "created_at", "updated_at"}).
			AddRow(album.ID, album.Title, album.ReleaseDate, album.CoverURL, album.CreatedAt, album.UpdatedAt)

		mock.ExpectQuery("SELECT (.+) FROM albums WHERE id = ?").
			WithArgs(albumID).
			WillReturnRows(rows)

		foundAlbum, err := repo.FindByID(albumID)
		assert.NoError(t, err)
		assert.Equal(t, album.ID, foundAlbum.ID)
		assert.Equal(t, album.Title, foundAlbum.Title)
		assert.Equal(t, album.CoverURL, foundAlbum.CoverURL)
	})

	// Сценарий с ошибкой
	t.Run("error", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM albums WHERE id = ?").
			WithArgs(albumID).
			WillReturnError(errors.New("db error"))

		foundAlbum, err := repo.FindByID(albumID)
		assert.Error(t, err)
		assert.Nil(t, foundAlbum)
	})

	// Проверка, что все ожидаемые запросы были выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAlbumRepository_Save(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewAlbumRepository(db)

	albumID := uuid.New()
	now := time.Now()
	album := &models.Album{
		ID:          albumID,
		Title:       "Test Album",
		ReleaseDate: now,
		CoverURL:    "http://example.com/cover.jpg",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Успешное сохранение
	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO albums").
			WithArgs(album.ID, album.Title, album.ReleaseDate, album.CoverURL, album.CreatedAt, album.UpdatedAt).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Save(album)
		assert.NoError(t, err)
	})

	// Ошибка при сохранении
	t.Run("error", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO albums").
			WithArgs(album.ID, album.Title, album.ReleaseDate, album.CoverURL, album.CreatedAt, album.UpdatedAt).
			WillReturnError(errors.New("db error"))

		err := repo.Save(album)
		assert.Error(t, err)
	})

	// Проверка, что все ожидаемые запросы были выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAlbumRepository_GetTracks(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewAlbumRepository(db)

	albumID := uuid.New()
	now := time.Now()
	tracks := []*models.Track{
		{
			ID:         uuid.New(),
			Title:      "Track 1",
			Duration:   180,
			AlbumID:    albumID,
			ArtistName: "Test Artist",
			AddedDate:  now,
			UpdatedAt:  now,
		},
		{
			ID:         uuid.New(),
			Title:      "Track 2",
			Duration:   240,
			AlbumID:    albumID,
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

		mock.ExpectQuery("SELECT (.+) FROM tracks WHERE album_id = ?").
			WithArgs(albumID).
			WillReturnRows(rows)

		foundTracks, err := repo.GetTracks(albumID)
		assert.NoError(t, err)
		assert.Len(t, foundTracks, 2)
		assert.Equal(t, tracks[0].Title, foundTracks[0].Title)
		assert.Equal(t, tracks[1].Title, foundTracks[1].Title)
	})

	// Ошибка при получении треков
	t.Run("error", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM tracks WHERE album_id = ?").
			WithArgs(albumID).
			WillReturnError(errors.New("db error"))

		foundTracks, err := repo.GetTracks(albumID)
		assert.Error(t, err)
		assert.Nil(t, foundTracks)
	})

	// Проверка, что все ожидаемые запросы были выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAlbumRepository_AddTrackToAlbum(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewAlbumRepository(db)

	albumID := uuid.New()
	trackID := uuid.New()

	// Успешное добавление трека в альбом
	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("UPDATE tracks SET album_id = \\$1 WHERE id = \\$2").
			WithArgs(albumID, trackID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.AddTrackToAlbum(albumID, trackID)
		assert.NoError(t, err)
	})

	// Ошибка при добавлении трека
	t.Run("error", func(t *testing.T) {
		mock.ExpectExec("UPDATE tracks SET album_id = \\$1 WHERE id = \\$2").
			WithArgs(albumID, trackID).
			WillReturnError(errors.New("db error"))

		err := repo.AddTrackToAlbum(albumID, trackID)
		assert.Error(t, err)
	})

	// Проверка, что все ожидаемые запросы были выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAlbumRepository_ListAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewAlbumRepository(db)

	now := time.Now()
	albums := []*models.Album{
		{
			ID:          uuid.New(),
			Title:       "Album 1",
			ReleaseDate: now,
			CoverURL:    "http://example.com/cover1.jpg",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          uuid.New(),
			Title:       "Album 2",
			ReleaseDate: now,
			CoverURL:    "http://example.com/cover2.jpg",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	// Успешное получение всех альбомов
	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title", "release_date", "cover_url", "created_at", "updated_at"})
		for _, album := range albums {
			rows.AddRow(album.ID, album.Title, album.ReleaseDate, album.CoverURL, album.CreatedAt, album.UpdatedAt)
		}

		mock.ExpectQuery("SELECT (.+) FROM albums").
			WillReturnRows(rows)

		foundAlbums, err := repo.ListAll()
		assert.NoError(t, err)
		assert.Len(t, foundAlbums, 2)
		assert.Equal(t, albums[0].Title, foundAlbums[0].Title)
		assert.Equal(t, albums[1].Title, foundAlbums[1].Title)
	})

	// Ошибка при получении альбомов
	t.Run("error", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM albums").
			WillReturnError(errors.New("db error"))

		foundAlbums, err := repo.ListAll()
		assert.Error(t, err)
		assert.Nil(t, foundAlbums)
	})

	// Проверка, что все ожидаемые запросы были выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
