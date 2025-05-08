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

func TestSessionRepository_CreateSession(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewSessionRepository(db)

	userID := uuid.New()

	// Успешное создание сессии
	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO sessions").
			WithArgs(sqlmock.AnyArg(), userID, sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		session, err := repo.CreateSession(userID)
		assert.NoError(t, err)
		assert.NotNil(t, session)
		assert.NotEmpty(t, session.ID)
		assert.Equal(t, userID, session.UserID)
		assert.True(t, session.ExpiresAt.After(time.Now()))
	})

	// Ошибка при создании сессии
	t.Run("error", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO sessions").
			WithArgs(sqlmock.AnyArg(), userID, sqlmock.AnyArg()).
			WillReturnError(errors.New("db error"))

		session, err := repo.CreateSession(userID)
		assert.Error(t, err)
		assert.Nil(t, session)
	})

	// Проверка, что все ожидаемые запросы были выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSessionRepository_GetSession(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewSessionRepository(db)

	sessionID := uuid.New().String()
	userID := uuid.New()
	expiresAt := time.Now().Add(24 * time.Hour)

	session := &models.Session{
		ID:        sessionID,
		UserID:    userID,
		ExpiresAt: expiresAt,
	}

	// Успешное получение сессии
	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "user_id", "expires_at"}).
			AddRow(session.ID, session.UserID, session.ExpiresAt)

		mock.ExpectQuery("SELECT id, user_id, expires_at FROM sessions WHERE id = \\$1 AND expires_at > NOW\\(\\)").
			WithArgs(sessionID).
			WillReturnRows(rows)

		foundSession, err := repo.GetSession(sessionID)
		assert.NoError(t, err)
		assert.Equal(t, session.ID, foundSession.ID)
		assert.Equal(t, session.UserID, foundSession.UserID)
	})

	// Сессия не найдена
	t.Run("not_found", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, user_id, expires_at FROM sessions WHERE id = \\$1 AND expires_at > NOW\\(\\)").
			WithArgs(sessionID).
			WillReturnError(errors.New("not found"))

		foundSession, err := repo.GetSession(sessionID)
		assert.Error(t, err)
		assert.Nil(t, foundSession)
	})

	// Проверка, что все ожидаемые запросы были выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSessionRepository_DeleteSession(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewSessionRepository(db)

	sessionID := uuid.New().String()

	// Успешное удаление сессии
	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM sessions WHERE id = \\$1").
			WithArgs(sessionID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.DeleteSession(sessionID)
		assert.NoError(t, err)
	})

	// Ошибка при удалении сессии
	t.Run("error", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM sessions WHERE id = \\$1").
			WithArgs(sessionID).
			WillReturnError(errors.New("db error"))

		err := repo.DeleteSession(sessionID)
		assert.Error(t, err)
	})

	// Проверка, что все ожидаемые запросы были выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSessionRepository_DeleteAllForUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewSessionRepository(db)

	userID := uuid.New()

	// Успешное удаление всех сессий пользователя
	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM sessions WHERE user_id = \\$1").
			WithArgs(userID).
			WillReturnResult(sqlmock.NewResult(0, 2))

		err := repo.DeleteAllForUser(userID)
		assert.NoError(t, err)
	})

	// Ошибка при удалении сессий
	t.Run("error", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM sessions WHERE user_id = \\$1").
			WithArgs(userID).
			WillReturnError(errors.New("db error"))

		err := repo.DeleteAllForUser(userID)
		assert.Error(t, err)
	})

	// Проверка, что все ожидаемые запросы были выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
