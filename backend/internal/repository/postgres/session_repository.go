package postgres

import (
	"database/sql"
	"fmt"
	"music-service/internal/models"
	"music-service/internal/repository/interfaces"
	"time"

	"github.com/google/uuid"
)

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) interfaces.SessionRepository {
	return &SessionRepository{
		db: db,
	}
}

func (r *SessionRepository) CreateSession(userID uuid.UUID, token string) (*models.Session, error) {
	session := &models.Session{
		ID:        uuid.New(),
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	query := `INSERT INTO sessions (id, user_id, token, expires_at) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(query, session.ID, userID, session.Token, session.ExpiresAt)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (r *SessionRepository) GetSession(sessionID string) (*models.Session, error) {
	var session models.Session
	query := `SELECT id, token, expires_at FROM sessions WHERE token = $1 AND expires_at > NOW()`
	err := r.db.QueryRow(query, sessionID).Scan(&session.ID, &session.Token, &session.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *SessionRepository) DeleteSession(sessionID string) error {
	query := `DELETE FROM sessions WHERE token = $1`
	_, err := r.db.Exec(query, sessionID)
	return err
}

func (r *SessionRepository) DeleteAllForUser(userID uuid.UUID) error {
	query := `DELETE FROM sessions WHERE user_id = $1`
	_, err := r.db.Exec(query, userID)
	return err
}

func (r *SessionRepository) GetSessionByToken(token string) (*models.Session, *models.User, error) {
	var session models.Session
	var user models.User
	var userID uuid.UUID

	fmt.Printf("Ищем сессию с токеном: %s\n", token)

	var count int
	countErr := r.db.QueryRow(`SELECT COUNT(*) FROM sessions WHERE token = $1`, token).Scan(&count)
	if countErr != nil {
		fmt.Printf("Ошибка при проверке наличия сессии: %v\n", countErr)
	} else {
		fmt.Printf("Найдено сессий: %d\n", count)
	}

	err := r.db.QueryRow(`
		SELECT s.id, s.user_id, s.expires_at, s.token, 
			   u.id, u.login, u.password, u.permission, u.created_at, u.updated_at
		FROM sessions s 
		JOIN users u ON s.user_id = u.id
		WHERE s.token = $1 AND s.expires_at > NOW()`,
		token).Scan(
		&session.ID, &userID, &session.ExpiresAt, &session.Token,
		&user.ID, &user.Login, &user.Password, &user.Permission, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		fmt.Printf("Ошибка при поиске сессии: %v\n", err)
		return nil, nil, err
	}

	fmt.Printf("Успешно найдена сессия для пользователя: %s\n", user.Login)
	return &session, &user, nil
}
