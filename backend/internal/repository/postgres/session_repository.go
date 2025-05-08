package postgres

import (
	"database/sql"
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

func (r *SessionRepository) CreateSession(userID uuid.UUID) (*models.Session, error) {
	session := &models.Session{
		ID:        uuid.New().String(),
		UserID:    userID,
		ExpiresAt: time.Now().Add(24 * time.Hour), // Сессия на 24 часа
	}

	query := `INSERT INTO sessions (id, user_id, expires_at) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(query, session.ID, session.UserID, session.ExpiresAt)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (r *SessionRepository) GetSession(sessionID string) (*models.Session, error) {
	var session models.Session
	query := `SELECT id, user_id, expires_at FROM sessions WHERE id = $1 AND expires_at > NOW()`
	err := r.db.QueryRow(query, sessionID).Scan(&session.ID, &session.UserID, &session.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *SessionRepository) DeleteSession(sessionID string) error {
	query := `DELETE FROM sessions WHERE id = $1`
	_, err := r.db.Exec(query, sessionID)
	return err
}

func (r *SessionRepository) DeleteAllForUser(userID uuid.UUID) error {
	query := `DELETE FROM sessions WHERE user_id = $1`
	_, err := r.db.Exec(query, userID)
	return err
}
