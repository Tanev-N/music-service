package handlers

import (
	"encoding/json"
	"fmt"
	"music-service/internal/models"
	"music-service/internal/usecases/interfaces"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type AlbumHandler struct {
	albumUseCase interfaces.AlbumUseCase
}

func NewAlbumHandler(albumUseCase interfaces.AlbumUseCase) *AlbumHandler {
	return &AlbumHandler{
		albumUseCase: albumUseCase,
	}
}

// Проверка прав администратора
func (h *AlbumHandler) isAdmin(r *http.Request) bool {
	// Для быстрого тестирования: проверяем непосредственно токен из Authorization
	authHeader := r.Header.Get("Authorization")

	// Если это наш фиксированный токен, то разрешаем
	if authHeader == "Bearer 33333333-3333-3333-3333-333333333333" {
		fmt.Println("Прямое разрешение по токену админа в AlbumHandler")
		return true
	}

	// Обычная проверка по заголовку X-User-Permission
	permission := r.Header.Get("X-User-Permission")
	isAdmin := permission == string(models.AdminPermission)
	fmt.Printf("AlbumHandler - Проверка прав админа: X-User-Permission=%s, isAdmin=%v\n", permission, isAdmin)
	return isAdmin
}

type createAlbumRequest struct {
	Title       string    `json:"title"`
	ReleaseDate time.Time `json:"release_date"`
	CoverURL    string    `json:"cover_url"`
	Artist      string    `json:"artist"`
}

type albumResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Artist      string    `json:"artist"`
	ReleaseDate time.Time `json:"release_date"`
	CoverURL    string    `json:"cover_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func toAlbumResponse(album *models.Album) albumResponse {
	return albumResponse{
		ID:          album.ID.String(),
		Title:       album.Title,
		Artist:      album.Artist,
		ReleaseDate: album.ReleaseDate,
		CoverURL:    album.CoverURL,
		CreatedAt:   album.CreatedAt,
		UpdatedAt:   album.UpdatedAt,
	}
}

func writeAlbumJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeAlbumError(w http.ResponseWriter, status int, message string) {
	writeAlbumJSON(w, status, errorResponse{Error: message})
}

// CreateAlbum создает новый альбом
func (h *AlbumHandler) CreateAlbum(w http.ResponseWriter, r *http.Request) {
	// Проверяем права администратора
	if !h.isAdmin(r) {
		writeAlbumError(w, http.StatusForbidden, "Доступ запрещен: требуются права администратора")
		return
	}

	var req createAlbumRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeAlbumError(w, http.StatusBadRequest, "Некорректные данные запроса")
		return
	}

	album, err := h.albumUseCase.CreateAlbum(req.Title, req.Artist, req.ReleaseDate, req.CoverURL)
	if err != nil {
		writeAlbumError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeAlbumJSON(w, http.StatusCreated, toAlbumResponse(album))
}

// GetAlbumDetails получает информацию об альбоме
func (h *AlbumHandler) GetAlbumDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	albumID, err := uuid.Parse(vars["id"])
	if err != nil {
		writeAlbumError(w, http.StatusBadRequest, "Недопустимый идентификатор альбома")
		return
	}

	album, tracks, err := h.albumUseCase.GetAlbumDetails(albumID)
	if err != nil {
		writeAlbumError(w, http.StatusNotFound, "Альбом не найден")
		return
	}

	type albumWithTracks struct {
		albumResponse
		Tracks []*models.Track `json:"tracks"`
	}

	response := albumWithTracks{
		albumResponse: toAlbumResponse(album),
		Tracks:        tracks,
	}

	writeAlbumJSON(w, http.StatusOK, response)
}

// UpdateAlbum обновляет информацию об альбоме
func (h *AlbumHandler) UpdateAlbum(w http.ResponseWriter, r *http.Request) {
	if !h.isAdmin(r) {
		writeAlbumError(w, http.StatusForbidden, "Доступ запрещен: требуются права администратора")
		return
	}

	vars := mux.Vars(r)
	albumID, err := uuid.Parse(vars["id"])
	if err != nil {
		writeAlbumError(w, http.StatusBadRequest, "Недопустимый идентификатор альбома")
		return
	}

	var req createAlbumRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeAlbumError(w, http.StatusBadRequest, "Некорректные данные запроса")
		return
	}

	if err := h.albumUseCase.UpdateAlbumInfo(albumID, req.Title, req.Artist, req.CoverURL, req.ReleaseDate); err != nil {
		writeAlbumError(w, http.StatusBadRequest, err.Error())
		return
	}

	album, _, err := h.albumUseCase.GetAlbumDetails(albumID)
	if err != nil {
		writeAlbumError(w, http.StatusInternalServerError, "Ошибка при получении обновленных данных альбома")
		return
	}

	writeAlbumJSON(w, http.StatusOK, toAlbumResponse(album))
}

// AddTrackToAlbum добавляет трек в альбом
func (h *AlbumHandler) AddTrackToAlbum(w http.ResponseWriter, r *http.Request) {
	if !h.isAdmin(r) {
		writeAlbumError(w, http.StatusForbidden, "Доступ запрещен: требуются права администратора")
		return
	}

	vars := mux.Vars(r)
	albumID, err := uuid.Parse(vars["id"])
	if err != nil {
		writeAlbumError(w, http.StatusBadRequest, "Недопустимый идентификатор альбома")
		return
	}

	var req struct {
		TrackID string `json:"track_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeAlbumError(w, http.StatusBadRequest, "Некорректные данные запроса")
		return
	}

	trackID, err := uuid.Parse(req.TrackID)
	if err != nil {
		writeAlbumError(w, http.StatusBadRequest, "Недопустимый идентификатор трека")
		return
	}

	if err := h.albumUseCase.AddTrackToAlbum(albumID, trackID); err != nil {
		writeAlbumError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

// RemoveTrackFromAlbum удаляет трек из альбома
func (h *AlbumHandler) RemoveTrackFromAlbum(w http.ResponseWriter, r *http.Request) {
	if !h.isAdmin(r) {
		writeAlbumError(w, http.StatusForbidden, "Доступ запрещен: требуются права администратора")
		return
	}

	vars := mux.Vars(r)
	albumID, err := uuid.Parse(vars["id"])
	if err != nil {
		writeAlbumError(w, http.StatusBadRequest, "Недопустимый идентификатор альбома")
		return
	}

	trackID, err := uuid.Parse(vars["track_id"])
	if err != nil {
		writeAlbumError(w, http.StatusBadRequest, "Недопустимый идентификатор трека")
		return
	}

	if err := h.albumUseCase.RemoveTrackFromAlbum(albumID, trackID); err != nil {
		writeAlbumError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

// ListAllAlbums выводит список всех альбомов
func (h *AlbumHandler) ListAllAlbums(w http.ResponseWriter, r *http.Request) {
	albums, err := h.albumUseCase.ListAll()
	if err != nil {
		writeAlbumError(w, http.StatusInternalServerError, "Ошибка при получении списка альбомов")
		return
	}

	var response []albumResponse
	for _, album := range albums {
		response = append(response, toAlbumResponse(album))
	}

	writeAlbumJSON(w, http.StatusOK, response)
}

// DeleteAlbum удаляет альбом
func (h *AlbumHandler) DeleteAlbum(w http.ResponseWriter, r *http.Request) {
	// Проверяем права администратора
	if !h.isAdmin(r) {
		writeAlbumError(w, http.StatusForbidden, "Доступ запрещен: требуются права администратора")
		return
	}

	vars := mux.Vars(r)
	albumID, err := uuid.Parse(vars["id"])
	if err != nil {
		writeAlbumError(w, http.StatusBadRequest, "Недопустимый идентификатор альбома")
		return
	}

	if err := h.albumUseCase.Delete(albumID); err != nil {
		writeAlbumError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
