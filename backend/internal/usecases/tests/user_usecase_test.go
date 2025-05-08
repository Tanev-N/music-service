package usecases_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"music-service/internal/models"
	"music-service/internal/repository/interfaces/mocks"
	"music-service/internal/usecases"

	"github.com/google/uuid"
)

func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	sessionRepo := mocks.NewMockSessionRepository(ctrl)
	uc := usecases.NewUserUseCase(userRepo, sessionRepo)

	expectedID := uuid.New()

	userRepo.EXPECT().Save(gomock.Any()).DoAndReturn(func(u *models.User) error {
		u.ID = expectedID
		return nil
	})

	tests := []struct {
		name          string
		login         string
		password      string
		mockSetup     func()
		expectedUser  *models.User
		expectedError string
	}{
		{
			name:     "successful registration",
			login:    "validlogin",
			password: "validpassword",
			mockSetup: func() {
				userRepo.EXPECT().Search("validlogin").Return([]*models.User{}, nil)
			},
			expectedUser: &models.User{
				ID:         expectedID,
				Login:      "validlogin",
				Password:   "hashed_validpassword",
				Permission: "user",
			},
		},
		{
			name:          "login too short",
			login:         "short",
			password:      "validpassword",
			mockSetup:     func() {},
			expectedError: "login must be at least 6 characters",
		},
		{
			name:          "password too short",
			login:         "validlogin",
			password:      "short",
			mockSetup:     func() {},
			expectedError: "password must be at least 8 characters",
		},
		{
			name:     "user already exists",
			login:    "existing",
			password: "validpassword",
			mockSetup: func() {
				userRepo.EXPECT().Search("existing").Return([]*models.User{
					{Login: "existing"},
				}, nil)
			},
			expectedError: "user already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			user, err := uc.Register(tt.login, tt.password)

			if tt.expectedError != "" {
				assert.ErrorContains(t, err, tt.expectedError)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser.ID, user.ID)
				assert.Equal(t, tt.expectedUser.Login, user.Login)
				assert.Equal(t, tt.expectedUser.Password, user.Password)
				assert.Equal(t, tt.expectedUser.Permission, user.Permission)
			}
		})
	}
}

func TestAuthenticate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	sessionRepo := mocks.NewMockSessionRepository(ctrl)
	uc := usecases.NewUserUseCase(userRepo, sessionRepo)
	tests := []struct {
		name            string
		login           string
		password        string
		mockSetup       func()
		expectedUser    *models.User
		expectedSession *models.Session
		expectedError   string
	}{
		{
			name:     "invalid credentials - user not found",
			login:    "invalid",
			password: "password",
			mockSetup: func() {
				userRepo.EXPECT().Search("invalid").Return([]*models.User{}, nil)
			},
			expectedError: "invalid credentials",
		},
		{
			name:     "invalid credentials - wrong password",
			login:    "validlogin",
			password: "wrongpassword",
			mockSetup: func() {
				userRepo.EXPECT().Search("validlogin").Return([]*models.User{
					{
						ID:       uuid.New(),
						Login:    "validlogin",
						Password: "hashed_validpassword",
					},
				}, nil)
			},
			expectedError: "invalid credentials",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			user, session, err := uc.Authenticate(tt.login, tt.password)

			if tt.expectedError != "" {
				assert.ErrorContains(t, err, tt.expectedError)
				assert.Nil(t, user)
				assert.Nil(t, session)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser.ID, user.ID)
				assert.Equal(t, tt.expectedUser.Login, user.Login)
				assert.Equal(t, tt.expectedSession.ID, session.ID)
				assert.Equal(t, tt.expectedSession.UserID, session.UserID)
			}
		})
	}
}

func TestGetUserProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	sessionRepo := mocks.NewMockSessionRepository(ctrl)
	uc := usecases.NewUserUseCase(userRepo, sessionRepo)

	t.Run("successful get", func(t *testing.T) {
		userID := uuid.New()
		userRepo.EXPECT().FindByID(userID).Return(&models.User{
			ID:    userID,
			Login: "testuser",
		}, nil)

		user, err := uc.GetUserProfile(userID)

		assert.NoError(t, err)
		assert.Equal(t, userID, user.ID)
		assert.Equal(t, "testuser", user.Login)
	})

	t.Run("user not found", func(t *testing.T) {
		userID := uuid.New()
		userRepo.EXPECT().FindByID(userID).Return(nil, errors.New("not found"))

		user, err := uc.GetUserProfile(userID)

		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestUpdatePermissions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	sessionRepo := mocks.NewMockSessionRepository(ctrl)
	uc := usecases.NewUserUseCase(userRepo, sessionRepo)

	// Фиксированный UUID для пользователя
	userID := uuid.MustParse("4645953f-d539-4fbf-bd2a-77316ab70998")

	tests := []struct {
		name        string
		userID      uuid.UUID
		permission  string
		mockSetup   func()
		expectedErr string
	}{
		{
			name:       "successful update",
			userID:     userID,
			permission: "admin",
			mockSetup: func() {
				userRepo.EXPECT().FindByID(userID).Return(&models.User{
					ID:         userID,
					Permission: "user",
				}, nil)
				userRepo.EXPECT().Save(gomock.Any()).Return(nil)
			},
		},
		{
			name:       "user not found",
			userID:     uuid.New(),
			permission: "admin",
			mockSetup: func() {
				userRepo.EXPECT().FindByID(gomock.Any()).Return(nil, errors.New("not found"))
			},
			expectedErr: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			err := uc.UpdatePermissions(tt.userID, models.Permission(tt.permission))

			if tt.expectedErr != "" {
				assert.ErrorContains(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mocks.NewMockUserRepository(ctrl)
	sessionRepo := mocks.NewMockSessionRepository(ctrl)
	uc := usecases.NewUserUseCase(userRepo, sessionRepo)

	userID := uuid.MustParse("4645953f-d539-4fbf-bd2a-77316ab70998")

	tests := []struct {
		name        string
		userID      uuid.UUID
		mockSetup   func()
		expectedErr string
	}{
		{
			name:   "delete failed",
			userID: userID,
			mockSetup: func() {
				sessionRepo.EXPECT().DeleteAllForUser(userID).Return(nil)
				userRepo.EXPECT().Delete(userID).Return(errors.New("delete failed"))
			},
			expectedErr: "delete failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			err := uc.DeleteUser(tt.userID)

			if tt.expectedErr != "" {
				assert.ErrorContains(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
