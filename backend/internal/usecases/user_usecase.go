package usecases

import (
	"errors"
	"music-service/internal/models"
	"music-service/internal/repository/interfaces"
	usecaseInterfaces "music-service/internal/usecases/interfaces"

	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, errors.New("error hashing password")
	}

	now := time.Now()
	user := &models.User{
		ID:         uuid.New(),
		Login:      login,
		Password:   hashedPassword,
		Permission: models.UserPermission,
		CreatedAt:  now,
		UpdatedAt:  now,
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

	session := &models.Session{
		ID:        uuid.New(),
		Token:     generateToken(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	createdSession, err := uc.sessionRepo.CreateSession(user.ID)
	if err != nil {
		return nil, nil, err
	}

	createdSession.Token = session.Token

	return user, createdSession, nil
}

func (uc *userUseCase) GetUserProfile(userID uuid.UUID) (*models.User, error) {
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (uc *userUseCase) UpdatePermissions(userID uuid.UUID, permission models.Permission) error {
	if !permission.IsValid() {
		return errors.New("invalid permission")
	}

	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	user.Permission = permission
	user.UpdatedAt = time.Now()

	return uc.userRepo.Save(user)
}

func (uc *userUseCase) DeleteUser(userID uuid.UUID) error {
	_, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	if err := uc.sessionRepo.DeleteAllForUser(userID); err != nil {
		return err
	}
	return uc.userRepo.Delete(userID)
}

func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateToken() string {
	return uuid.New().String()
}
