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

func TestUserRepository_FindByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewUserRepository(db)

	userID := uuid.New()
	user := &models.User{
		ID:         userID,
		Login:      "testuser",
		Password:   "hashedpassword",
		Permission: models.UserPermission,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Успешный сценарий
	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "login", "password", "permission", "created_at", "updated_at"}).
			AddRow(user.ID, user.Login, user.Password, user.Permission, user.CreatedAt, user.UpdatedAt)

		mock.ExpectQuery("SELECT (.+) FROM users WHERE id = \\$1").
			WithArgs(userID).
			WillReturnRows(rows)

		foundUser, err := repo.FindByID(userID)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, foundUser.ID)
		assert.Equal(t, user.Login, foundUser.Login)
		assert.Equal(t, user.Password, foundUser.Password)
		assert.Equal(t, user.Permission, foundUser.Permission)
	})

	// Сценарий с ошибкой
	t.Run("error", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM users WHERE id = \\$1").
			WithArgs(userID).
			WillReturnError(errors.New("db error"))

		foundUser, err := repo.FindByID(userID)
		assert.Error(t, err)
		assert.Nil(t, foundUser)
	})

	// Проверка, что все ожидаемые запросы были выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepository_Save(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewUserRepository(db)

	userID := uuid.New()
	user := &models.User{
		ID:         userID,
		Login:      "testuser",
		Password:   "hashedpassword",
		Permission: models.UserPermission,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Успешное сохранение
	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO users").
			WithArgs(user.ID, user.Login, user.Password, user.Permission, user.CreatedAt, user.UpdatedAt).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Save(user)
		assert.NoError(t, err)
	})

	// Ошибка при сохранении
	t.Run("error", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO users").
			WithArgs(user.ID, user.Login, user.Password, user.Permission, user.CreatedAt, user.UpdatedAt).
			WillReturnError(errors.New("db error"))

		err := repo.Save(user)
		assert.Error(t, err)
	})

	// Проверка, что все ожидаемые запросы были выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewUserRepository(db)

	userID := uuid.New()

	// Успешное удаление
	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM users WHERE id = \\$1").
			WithArgs(userID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Delete(userID)
		assert.NoError(t, err)
	})

	// Ошибка при удалении
	t.Run("error", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM users WHERE id = \\$1").
			WithArgs(userID).
			WillReturnError(errors.New("db error"))

		err := repo.Delete(userID)
		assert.Error(t, err)
	})

	// Проверка, что все ожидаемые запросы были выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUserRepository_Search(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := postgres.NewUserRepository(db)

	searchQuery := "test"
	users := []*models.User{
		{
			ID:         uuid.New(),
			Login:      "testuser1",
			Password:   "hashedpassword1",
			Permission: models.UserPermission,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			ID:         uuid.New(),
			Login:      "testuser2",
			Password:   "hashedpassword2",
			Permission: models.UserPermission,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	// Успешный поиск
	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "login", "password", "permission", "created_at", "updated_at"})
		for _, user := range users {
			rows.AddRow(user.ID, user.Login, user.Password, user.Permission, user.CreatedAt, user.UpdatedAt)
		}

		mock.ExpectQuery("SELECT (.+) FROM users WHERE login ILIKE \\$1").
			WithArgs("%" + searchQuery + "%").
			WillReturnRows(rows)

		foundUsers, err := repo.Search(searchQuery)
		assert.NoError(t, err)
		assert.Len(t, foundUsers, 2)
		assert.Equal(t, users[0].Login, foundUsers[0].Login)
		assert.Equal(t, users[1].Login, foundUsers[1].Login)
	})

	// Ошибка при поиске
	t.Run("error", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM users WHERE login ILIKE \\$1").
			WithArgs("%" + searchQuery + "%").
			WillReturnError(errors.New("db error"))

		foundUsers, err := repo.Search(searchQuery)
		assert.Error(t, err)
		assert.Nil(t, foundUsers)
	})

	// Пустой результат
	t.Run("empty_result", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "login", "password", "permission", "created_at", "updated_at"})

		mock.ExpectQuery("SELECT (.+) FROM users WHERE login ILIKE \\$1").
			WithArgs("%" + searchQuery + "%").
			WillReturnRows(rows)

		foundUsers, err := repo.Search(searchQuery)
		assert.NoError(t, err)
		assert.Empty(t, foundUsers)
	})

	// Проверка, что все ожидаемые запросы были выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
