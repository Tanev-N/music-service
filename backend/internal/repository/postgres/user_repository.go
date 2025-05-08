package postgres

import (
	"database/sql"
	"music-service/internal/models"
	"music-service/internal/repository/interfaces"

	"github.com/google/uuid"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) interfaces.UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	var permissionStr string

	query := `SELECT id, login, password, permission, created_at, updated_at FROM users WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Login,
		&user.Password,
		&permissionStr,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	user.Permission = models.Permission(permissionStr)
	return &user, nil
}

func (r *UserRepository) Save(user *models.User) error {
	query := `
		INSERT INTO users (id, login, password, permission, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE 
		SET login = $2, password = $3, permission = $4, updated_at = $6
	`
	_, err := r.db.Exec(query, user.ID, user.Login, user.Password, user.Permission, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *UserRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *UserRepository) Search(query string) ([]*models.User, error) {
	var users []*models.User
	rows, err := r.db.Query(`SELECT id, login, password, permission, created_at, updated_at FROM users WHERE login ILIKE $1`, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		var permissionStr string

		err := rows.Scan(
			&user.ID,
			&user.Login,
			&user.Password,
			&permissionStr,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		user.Permission = models.Permission(permissionStr)
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
