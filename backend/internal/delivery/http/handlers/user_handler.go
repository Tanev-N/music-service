package handlers

import (
	"encoding/json"
	"music-service/internal/models"
	"music-service/internal/usecases/interfaces"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	userUseCase interfaces.UserUseCase
}

func NewUserHandler(userUseCase interfaces.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

type registerRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type authRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type authResponse struct {
	User    *models.User    `json:"user"`
	Session *models.Session `json:"session"`
}

type updatePermissionsRequest struct {
	Permission string `json:"permission"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message})
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Некорректные данные запроса")
		return
	}

	user, err := h.userUseCase.Register(req.Login, req.Password)
	if err != nil {
		switch err.Error() {
		case "user already exists":
			writeError(w, http.StatusBadRequest, "Пользователь уже существует")
		case "login must be at least 6 characters", "password must be at least 8 characters":
			writeError(w, http.StatusBadRequest, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "Ошибка сервера")
		}
		return
	}

	writeJSON(w, http.StatusCreated, user)
}

func (h *UserHandler) AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Некорректные данные запроса")
		return
	}

	user, session, err := h.userUseCase.Authenticate(req.Login, req.Password)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Неверные учетные данные")
		return
	}

	writeJSON(w, http.StatusOK, authResponse{
		User:    user,
		Session: session,
	})
}

func (h *UserHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := uuid.Parse(vars["id"])
	if err != nil {
		writeError(w, http.StatusBadRequest, "Неверный формат ID")
		return
	}

	user, err := h.userUseCase.GetUserProfile(userID)
	if err != nil {
		writeError(w, http.StatusNotFound, "Пользователь не найден")
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func (h *UserHandler) UpdateUserPermissions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := uuid.Parse(vars["id"])
	if err != nil {
		writeError(w, http.StatusBadRequest, "Неверный формат ID")
		return
	}

	var req updatePermissionsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Некорректные данные запроса")
		return
	}

	permission := models.Permission(req.Permission)
	if !permission.IsValid() {
		writeError(w, http.StatusBadRequest, "Недопустимые права доступа")
		return
	}

	if err := h.userUseCase.UpdatePermissions(userID, permission); err != nil {
		if err.Error() == "user not found" {
			writeError(w, http.StatusNotFound, "Пользователь не найден")
			return
		}
		writeError(w, http.StatusInternalServerError, "Ошибка сервера")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := uuid.Parse(vars["id"])
	if err != nil {
		writeError(w, http.StatusBadRequest, "Неверный формат ID")
		return
	}

	if err := h.userUseCase.DeleteUser(userID); err != nil {
		if err.Error() == "user not found" {
			writeError(w, http.StatusNotFound, "Пользователь не найден")
			return
		}
		writeError(w, http.StatusInternalServerError, "Ошибка сервера")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
