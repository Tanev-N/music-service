package postgres

import (
	"database/sql"
	"music-service/internal/models"
	"music-service/internal/repository/interfaces"
	"time"

	"github.com/google/uuid"
)

type HistoryRepository struct {
	db *sql.DB
}

func NewHistoryRepository(db *sql.DB) interfaces.HistoryRepository {
	return &HistoryRepository{
		db: db,
	}
}

func (r *HistoryRepository) AddEntry(userID uuid.UUID, trackID uuid.UUID) error {
	query := `
		INSERT INTO listening_history (user_id, track_id, listened_at) 
		VALUES ($1, $2, $3)
	`
	_, err := r.db.Exec(query, userID, trackID, time.Now())
	return err
}

func (r *HistoryRepository) GetHistory(userID uuid.UUID) ([]*models.ListeningHistory, error) {
	var history []*models.ListeningHistory
	query := `
		SELECT lh.id, lh.user_id, lh.track_id, lh.listened_at
		FROM listening_history lh
		WHERE lh.user_id = $1
		ORDER BY lh.listened_at DESC
		LIMIT 100
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var entry models.ListeningHistory
		err := rows.Scan(
			&entry.ID,
			&entry.UserID,
			&entry.TrackID,
			&entry.ListenedAt,
		)
		if err != nil {
			return nil, err
		}
		history = append(history, &entry)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return history, nil
}

// GetPlayCount возвращает количество прослушиваний трека
func (r *HistoryRepository) GetPlayCount(trackID uuid.UUID) (int, error) {
	var count int
	query := `
		SELECT COUNT(*) 
		FROM listening_history 
		WHERE track_id = $1
	`
	err := r.db.QueryRow(query, trackID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
