package tests

import (
	"errors"
	"music-service/internal/models"
	"music-service/internal/repository/postgres"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGenreRepository_FindByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewGenreRepository(db)

	genreID := uuid.New()
	genre := &models.Genre{
		ID:   genreID,
		Name: "Test Genre",
	}

	// Успешный сценарий
	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name"}).
			AddRow(genre.ID, genre.Name)

		mock.ExpectQuery("SELECT id, name FROM genres WHERE id = ?").
			WithArgs(genreID).
			WillReturnRows(rows)

		foundGenre, err := repo.FindByID(genreID)
		assert.NoError(t, err)
		assert.Equal(t, genre.ID, foundGenre.ID)
		assert.Equal(t, genre.Name, foundGenre.Name)
	})

	// Сценарий с ошибкой
	t.Run("error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, name FROM genres WHERE id = ?").
			WithArgs(genreID).
			WillReturnError(errors.New("db error"))

		foundGenre, err := repo.FindByID(genreID)
		assert.Error(t, err)
		assert.Nil(t, foundGenre)
	})

	// Проверка, что все ожидаемые запросы были выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGenreRepository_Save(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewGenreRepository(db)

	genreID := uuid.New()
	genre := &models.Genre{
		ID:   genreID,
		Name: "Test Genre",
	}

	// Успешное сохранение
	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO genres").
			WithArgs(genre.ID, genre.Name).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Save(genre)
		assert.NoError(t, err)
	})

	// Ошибка при сохранении
	t.Run("error", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO genres").
			WithArgs(genre.ID, genre.Name).
			WillReturnError(errors.New("db error"))

		err := repo.Save(genre)
		assert.Error(t, err)
	})

	// Проверка, что все ожидаемые запросы были выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGenreRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewGenreRepository(db)

	genreID := uuid.New()

	// Успешное удаление
	t.Run("success", func(t *testing.T) {
		// Сначала удаляются связи с треками
		mock.ExpectExec("DELETE FROM track_genres WHERE genre_id = ?").
			WithArgs(genreID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// Затем удаляется сам жанр
		mock.ExpectExec("DELETE FROM genres WHERE id = ?").
			WithArgs(genreID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Delete(genreID)
		assert.NoError(t, err)
	})

	// Ошибка при удалении связей
	t.Run("error_track_genres", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM track_genres WHERE genre_id = ?").
			WithArgs(genreID).
			WillReturnError(errors.New("db error"))

		err := repo.Delete(genreID)
		assert.Error(t, err)
	})

	// Ошибка при удалении жанра
	t.Run("error_genre", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM track_genres WHERE genre_id = ?").
			WithArgs(genreID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectExec("DELETE FROM genres WHERE id = ?").
			WithArgs(genreID).
			WillReturnError(errors.New("db error"))

		err := repo.Delete(genreID)
		assert.Error(t, err)
	})

	// Проверка, что все ожидаемые запросы были выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGenreRepository_GetGenresForTrack(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewGenreRepository(db)

	trackID := uuid.New()
	genres := []*models.Genre{
		{
			ID:   uuid.New(),
			Name: "Rock",
		},
		{
			ID:   uuid.New(),
			Name: "Pop",
		},
	}

	// Успешное получение жанров
	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name"})
		for _, genre := range genres {
			rows.AddRow(genre.ID, genre.Name)
		}

		mock.ExpectQuery("SELECT g.id, g.name FROM genres g JOIN track_genres tg ON g.id = tg.genre_id WHERE tg.track_id = ?").
			WithArgs(trackID).
			WillReturnRows(rows)

		foundGenres, err := repo.GetGenresForTrack(trackID)
		assert.NoError(t, err)
		assert.Len(t, foundGenres, 2)
		assert.Equal(t, genres[0].Name, foundGenres[0].Name)
		assert.Equal(t, genres[1].Name, foundGenres[1].Name)
	})

	// Ошибка при получении жанров
	t.Run("error", func(t *testing.T) {
		mock.ExpectQuery("SELECT g.id, g.name FROM genres g JOIN track_genres tg ON g.id = tg.genre_id WHERE tg.track_id = ?").
			WithArgs(trackID).
			WillReturnError(errors.New("db error"))

		foundGenres, err := repo.GetGenresForTrack(trackID)
		assert.Error(t, err)
		assert.Nil(t, foundGenres)
	})

	// Проверка, что все ожидаемые запросы были выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGenreRepository_AddGenreToTrack(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewGenreRepository(db)

	trackID := uuid.New()
	genreID := uuid.New()

	// Успешное добавление жанра к треку
	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO track_genres").
			WithArgs(trackID, genreID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.AddGenreToTrack(trackID, genreID)
		assert.NoError(t, err)
	})

	// Ошибка при добавлении жанра
	t.Run("error", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO track_genres").
			WithArgs(trackID, genreID).
			WillReturnError(errors.New("db error"))

		err := repo.AddGenreToTrack(trackID, genreID)
		assert.Error(t, err)
	})

	// Проверка, что все ожидаемые запросы были выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGenreRepository_RemoveGenreFromTrack(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewGenreRepository(db)

	trackID := uuid.New()
	genreID := uuid.New()

	// Успешное удаление жанра из трека
	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM track_genres WHERE track_id = \\$1 AND genre_id = \\$2").
			WithArgs(trackID, genreID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.RemoveGenreFromTrack(trackID, genreID)
		assert.NoError(t, err)
	})

	// Ошибка при удалении жанра
	t.Run("error", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM track_genres WHERE track_id = \\$1 AND genre_id = \\$2").
			WithArgs(trackID, genreID).
			WillReturnError(errors.New("db error"))

		err := repo.RemoveGenreFromTrack(trackID, genreID)
		assert.Error(t, err)
	})

	// Проверка, что все ожидаемые запросы были выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGenreRepository_ListAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewGenreRepository(db)

	genres := []*models.Genre{
		{
			ID:   uuid.New(),
			Name: "Rock",
		},
		{
			ID:   uuid.New(),
			Name: "Pop",
		},
		{
			ID:   uuid.New(),
			Name: "Jazz",
		},
	}

	// Успешное получение всех жанров
	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name"})
		for _, genre := range genres {
			rows.AddRow(genre.ID, genre.Name)
		}

		mock.ExpectQuery("SELECT id, name FROM genres").
			WillReturnRows(rows)

		foundGenres, err := repo.ListAll()
		assert.NoError(t, err)
		assert.Len(t, foundGenres, 3)
		assert.Equal(t, genres[0].Name, foundGenres[0].Name)
		assert.Equal(t, genres[1].Name, foundGenres[1].Name)
		assert.Equal(t, genres[2].Name, foundGenres[2].Name)
	})

	// Ошибка при получении жанров
	t.Run("error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, name FROM genres").
			WillReturnError(errors.New("db error"))

		foundGenres, err := repo.ListAll()
		assert.Error(t, err)
		assert.Nil(t, foundGenres)
	})

	// Проверка, что все ожидаемые запросы были выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
