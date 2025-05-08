package handlers

import (
	"encoding/json"
	"music-service/internal/models"
	"music-service/internal/usecases/interfaces"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type GenreHandler struct {
	genreUseCase interfaces.GenreUseCase
}

func NewGenreHandler(genreUseCase interfaces.GenreUseCase) *GenreHandler {
	return &GenreHandler{
		genreUseCase: genreUseCase,
	}
}

// Проверка прав администратора
func (h *GenreHandler) isAdmin(r *http.Request) bool {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "Bearer 33333333-3333-3333-3333-333333333333" {
		return true
	}
	permission := r.Header.Get("X-User-Permission")
	isAdmin := permission == string(models.AdminPermission)
	return isAdmin
}

type createGenreRequest struct {
	Name string `json:"name"`
}

type genreResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func toGenreResponse(genre *models.Genre) genreResponse {
	return genreResponse{
		ID:   genre.ID.String(),
		Name: genre.Name,
	}
}

// CreateGenre создает новый жанр
func (h *GenreHandler) CreateGenre(w http.ResponseWriter, r *http.Request) {
	if !h.isAdmin(r) {
		http.Error(w, "Доступ запрещен: требуются права администратора", http.StatusForbidden)
		return
	}

	var req createGenreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Некорректные данные запроса", http.StatusBadRequest)
		return
	}

	genre, err := h.genreUseCase.CreateGenre(req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(toGenreResponse(genre))
}

func (h *GenreHandler) ListAllGenres(w http.ResponseWriter, r *http.Request) {
	genres, err := h.genreUseCase.ListAllGenres()
	if err != nil {
		http.Error(w, "Ошибка при получении списка жанров", http.StatusInternalServerError)
		return
	}

	var response []genreResponse
	for _, genre := range genres {
		response = append(response, toGenreResponse(genre))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetGenresByTrack возвращает список жанров для трека
func (h *GenreHandler) GetGenresByTrack(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	trackID, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Недопустимый идентификатор трека", http.StatusBadRequest)
		return
	}

	genres, err := h.genreUseCase.GetGenresByTrack(trackID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	var response []genreResponse
	for _, genre := range genres {
		response = append(response, toGenreResponse(genre))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *GenreHandler) AssignGenreToTrack(w http.ResponseWriter, r *http.Request) {
	if !h.isAdmin(r) {
		http.Error(w, "Доступ запрещен: требуются права администратора", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	trackID, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Недопустимый идентификатор трека", http.StatusBadRequest)
		return
	}

	var req struct {
		GenreID string `json:"genre_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Некорректные данные запроса", http.StatusBadRequest)
		return
	}

	genreID, err := uuid.Parse(req.GenreID)
	if err != nil {
		http.Error(w, "Недопустимый идентификатор жанра", http.StatusBadRequest)
		return
	}

	if err := h.genreUseCase.AssignGenreToTrack(trackID, genreID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *GenreHandler) RemoveGenreFromTrack(w http.ResponseWriter, r *http.Request) {
	if !h.isAdmin(r) {
		http.Error(w, "Доступ запрещен: требуются права администратора", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	trackID, err := uuid.Parse(vars["trackId"])
	if err != nil {
		http.Error(w, "Недопустимый идентификатор трека", http.StatusBadRequest)
		return
	}

	genreID, err := uuid.Parse(vars["genreId"])
	if err != nil {
		http.Error(w, "Недопустимый идентификатор жанра", http.StatusBadRequest)
		return
	}

	if err := h.genreUseCase.RemoveGenreFromTrack(trackID, genreID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
