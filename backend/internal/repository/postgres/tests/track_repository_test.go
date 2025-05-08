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

func TestTrackRepository_FindByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewTrackRepository(db)

	trackID := uuid.New()
	albumID := uuid.New()
	now := time.Now()
	track := &models.Track{
		ID:         trackID,
		Title:      "Test Track",
		Duration:   240,
		FilePath:   "/tracks/test.mp3",
		AlbumID:    albumID,
		ArtistName: "Test Artist",
		CoverURL:   "http://example.com/cover.jpg",
		AddedDate:  now,
		UpdatedAt:  now,
		PlayCount:  5,
	}

	// Успешный сценарий
	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title", "duration", "file_path", "album_id", "artist_name", "cover_url", "added_date", "updated_at", "play_count"}).
			AddRow(track.ID, track.Title, track.Duration, track.FilePath, track.AlbumID, track.ArtistName, track.CoverURL, track.AddedDate, track.UpdatedAt, track.PlayCount)

		mock.ExpectQuery("SELECT (.+) FROM tracks WHERE id = ?").
			WithArgs(trackID).
			WillReturnRows(rows)

		foundTrack, err := repo.FindByID(trackID)
		assert.NoError(t, err)
		assert.Equal(t, track.ID, foundTrack.ID)
		assert.Equal(t, track.Title, foundTrack.Title)
		assert.Equal(t, track.Duration, foundTrack.Duration)
		assert.Equal(t, track.AlbumID, foundTrack.AlbumID)
		assert.Equal(t, track.ArtistName, foundTrack.ArtistName)
	})

	// Сценарий с ошибкой
	t.Run("error", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM tracks WHERE id = ?").
			WithArgs(trackID).
			WillReturnError(errors.New("db error"))

		foundTrack, err := repo.FindByID(trackID)
		assert.Error(t, err)
		assert.Nil(t, foundTrack)
	})

	// Проверка, что все ожидаемые запросы были выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTrackRepository_Save(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewTrackRepository(db)

	trackID := uuid.New()
	albumID := uuid.New()
	now := time.Now()
	track := &models.Track{
		ID:         trackID,
		Title:      "Test Track",
		Duration:   240,
		FilePath:   "/tracks/test.mp3",
		AlbumID:    albumID,
		ArtistName: "Test Artist",
		CoverURL:   "http://example.com/cover.jpg",
		AddedDate:  now,
		UpdatedAt:  now,
		PlayCount:  5,
	}

	// Успешное сохранение
	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO tracks").
			WithArgs(track.ID, track.Title, track.Duration, track.FilePath, track.AlbumID, track.ArtistName, track.CoverURL, track.AddedDate, track.UpdatedAt, track.PlayCount).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Save(track)
		assert.NoError(t, err)
	})

	// Ошибка при сохранении
	t.Run("error", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO tracks").
			WithArgs(track.ID, track.Title, track.Duration, track.FilePath, track.AlbumID, track.ArtistName, track.CoverURL, track.AddedDate, track.UpdatedAt, track.PlayCount).
			WillReturnError(errors.New("db error"))

		err := repo.Save(track)
		assert.Error(t, err)
	})

	// Проверка, что все ожидаемые запросы были выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTrackRepository_IncrementPlayCount(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewTrackRepository(db)

	trackID := uuid.New()

	// Успешное обновление счетчика
	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("UPDATE tracks SET play_count = play_count \\+ 1 WHERE id = ?").
			WithArgs(trackID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.IncrementPlayCount(trackID)
		assert.NoError(t, err)
	})

	// Ошибка при обновлении
	t.Run("error", func(t *testing.T) {
		mock.ExpectExec("UPDATE tracks SET play_count = play_count \\+ 1 WHERE id = ?").
			WithArgs(trackID).
			WillReturnError(errors.New("db error"))

		err := repo.IncrementPlayCount(trackID)
		assert.Error(t, err)
	})

	// Проверка, что все ожидаемые запросы были выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
