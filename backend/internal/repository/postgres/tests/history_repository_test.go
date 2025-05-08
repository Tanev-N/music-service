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

func TestHistoryRepository_AddEntry(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewHistoryRepository(db)

	userID := uuid.New()
	trackID := uuid.New()

	// Успешное добавление записи в историю прослушивания
	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO listening_history").
			WithArgs(userID, trackID, sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.AddEntry(userID, trackID)
		assert.NoError(t, err)
	})

	// Ошибка при добавлении записи
	t.Run("error", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO listening_history").
			WithArgs(userID, trackID, sqlmock.AnyArg()).
			WillReturnError(errors.New("db error"))

		err := repo.AddEntry(userID, trackID)
		assert.Error(t, err)
	})

	// Проверка, что все ожидаемые запросы были выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestHistoryRepository_GetHistory(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewHistoryRepository(db)

	userID := uuid.New()
	now := time.Now()

	historyEntries := []*models.ListeningHistory{
		{
			ID:         uuid.New(),
			UserID:     userID,
			TrackID:    uuid.New(),
			ListenedAt: now.Add(-1 * time.Hour),
		},
		{
			ID:         uuid.New(),
			UserID:     userID,
			TrackID:    uuid.New(),
			ListenedAt: now.Add(-2 * time.Hour),
		},
	}

	// Успешное получение истории прослушивания
	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "user_id", "track_id", "listened_at"})
		for _, entry := range historyEntries {
			rows.AddRow(entry.ID, entry.UserID, entry.TrackID, entry.ListenedAt)
		}

		mock.ExpectQuery("SELECT lh.id, lh.user_id, lh.track_id, lh.listened_at FROM listening_history lh WHERE lh.user_id = \\$1").
			WithArgs(userID).
			WillReturnRows(rows)

		history, err := repo.GetHistory(userID)
		assert.NoError(t, err)
		assert.Len(t, history, 2)
		assert.Equal(t, historyEntries[0].TrackID, history[0].TrackID)
		assert.Equal(t, historyEntries[1].TrackID, history[1].TrackID)
	})

	// Пустая история
	t.Run("empty_history", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "user_id", "track_id", "listened_at"})

		mock.ExpectQuery("SELECT lh.id, lh.user_id, lh.track_id, lh.listened_at FROM listening_history lh WHERE lh.user_id = \\$1").
			WithArgs(userID).
			WillReturnRows(rows)

		history, err := repo.GetHistory(userID)
		assert.NoError(t, err)
		assert.Empty(t, history)
	})

	// Ошибка при получении истории
	t.Run("error", func(t *testing.T) {
		mock.ExpectQuery("SELECT lh.id, lh.user_id, lh.track_id, lh.listened_at FROM listening_history lh WHERE lh.user_id = \\$1").
			WithArgs(userID).
			WillReturnError(errors.New("db error"))

		history, err := repo.GetHistory(userID)
		assert.Error(t, err)
		assert.Nil(t, history)
	})

	// Проверка, что все ожидаемые запросы были выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
