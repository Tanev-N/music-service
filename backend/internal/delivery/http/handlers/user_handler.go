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

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	user, err := h.userUseCase.Register(req.Login, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	user, session, err := h.userUseCase.Authenticate(req.Login, req.Password)
	if err != nil {
		http.Error(w, "Неверные учетные данные", http.StatusUnauthorized)
		return
	}

	response := authResponse{
		User:    user,
		Session: session,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	user, err := h.userUseCase.GetUserProfile(userID)
	if err != nil {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) UpdateUserPermissions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	var req updatePermissionsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	permission := models.Permission(req.Permission)
	if permission != models.AdminPermission && permission != models.UserPermission {
		http.Error(w, "Недопустимые права доступа", http.StatusBadRequest)
		return
	}

	err = h.userUseCase.UpdatePermissions(userID, permission)
	if err != nil {
		http.Error(w, "Ошибка обновления прав", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	err = h.userUseCase.DeleteUser(userID)
	if err != nil {
		http.Error(w, "Ошибка удаления пользователя", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type ErrorResponse struct {
	Error string `json:"error"`
}
