package usecases

import (
	"errors"
	"music-service/internal/models"
	"music-service/internal/repository/interfaces"
	usecaseInterfaces "music-service/internal/usecases/interfaces"

	"github.com/google/uuid"
)

type userUseCase struct {
	userRepo    interfaces.UserRepository
	sessionRepo interfaces.SessionRepository
}

func NewUserUseCase(
	userRepo interfaces.UserRepository,
	sessionRepo interfaces.SessionRepository,
) usecaseInterfaces.UserUseCase {
	return &userUseCase{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

func (uc *userUseCase) Register(login, password string) (*models.User, error) {
	if len(login) < 6 {
		return nil, errors.New("login must be at least 6 characters")
	}
	if len(password) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}

	users, err := uc.userRepo.Search(login)
	if err != nil {
		return nil, err
	}
	if len(users) > 0 {
		return nil, errors.New("user already exists")
	}

	user := &models.User{
		Login:      login,
		Password:   hashPassword(password),
		Permission: models.UserPermission,
	}

	if err := uc.userRepo.Save(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *userUseCase) Authenticate(login, password string) (*models.User, *models.Session, error) {
	users, err := uc.userRepo.Search(login)
	if err != nil || len(users) == 0 {
		return nil, nil, errors.New("invalid credentials")
	}

	user := users[0]
	if !checkPasswordHash(password, user.Password) {
		return nil, nil, errors.New("invalid credentials")
	}

	session, err := uc.sessionRepo.CreateSession(user.ID)
	if err != nil {
		return nil, nil, err
	}

	return user, session, nil
}

func (uc *userUseCase) GetUserProfile(userID uuid.UUID) (*models.User, error) {
	return uc.userRepo.FindByID(userID)
}

func (uc *userUseCase) UpdatePermissions(userID uuid.UUID, permission models.Permission) error {
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	user.Permission = permission
	return uc.userRepo.Save(user)
}

func (uc *userUseCase) DeleteUser(userID uuid.UUID) error {
	if err := uc.sessionRepo.DeleteAllForUser(userID); err != nil {
		return err
	}
	return uc.userRepo.Delete(userID)
}

func hashPassword(password string) string {
	return "hashed_" + password
}

func checkPasswordHash(password, hash string) bool {
	return hash == "hashed_"+password
}
