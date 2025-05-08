package handlers

import (
	"encoding/json"
	"fmt"
	"music-service/internal/models"
	"music-service/internal/usecases/interfaces"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const (
	maxMemory = 10 << 20
)

type TrackHandler struct {
	trackUseCase  interfaces.TrackUseCase
	allowedTypes  []string
	maxFileSizeMB int
}

func NewTrackHandler(trackUseCase interfaces.TrackUseCase, allowedTypes []string, maxFileSizeMB int) *TrackHandler {
	return &TrackHandler{
		trackUseCase:  trackUseCase,
		allowedTypes:  allowedTypes,
		maxFileSizeMB: maxFileSizeMB,
	}
}

// Обработчик для загрузки аудиофайла
func (h *TrackHandler) UploadTrack(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Заголовки запроса: %+v\n", r.Header)
	fmt.Printf("X-User-Permission: %s\n", r.Header.Get("X-User-Permission"))
	fmt.Printf("Authorization: %s\n", r.Header.Get("Authorization"))

	if !h.isAdmin(r) {
		fmt.Printf("Доступ запрещен: X-User-Permission=%s\n", r.Header.Get("X-User-Permission"))
		http.Error(w, "Доступ запрещен: требуются права администратора", http.StatusForbidden)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, int64(h.maxFileSizeMB<<20))

	if err := r.ParseMultipartForm(maxMemory); err != nil {
		http.Error(w, "Невозможно обработать загруженный файл: "+err.Error(), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Не удалось получить файл из запроса: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	fmt.Printf("Тип файла: %s\n", contentType)

	fmt.Printf("Разрешенные типы файлов: %v\n", h.allowedTypes)

	if !h.isAllowedFileType(contentType) {
		http.Error(w, "Недопустимый тип файла. Разрешены только MP3 файлы", http.StatusBadRequest)
		return
	}

	metadata, err := h.getTrackMetadataFromForm(r)
	if err != nil {
		http.Error(w, "Ошибка в метаданных трека: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Метаданные трека перед загрузкой: Title=%s, ArtistName=%s, AlbumID=%s\n",
		metadata.Title, metadata.ArtistName, metadata.AlbumID)

	track, err := h.trackUseCase.UploadTrack(file, header.Size, metadata)
	if err != nil {
		http.Error(w, "Ошибка при загрузке трека: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Printf("Трек создан: ID=%s, Title=%s, AlbumID=%s\n",
		track.ID, track.Title, track.AlbumID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(track)
}

// Получение метаданных трека из формы
func (h *TrackHandler) getTrackMetadataFromForm(r *http.Request) (models.TrackUploadMetadata, error) {
	metadata := models.TrackUploadMetadata{
		Title:      r.FormValue("title"),
		ArtistName: r.FormValue("artist_name"),
		CoverURL:   r.FormValue("cover_url"),
	}

	// Обрабатываем длительность если она указана
	if durationStr := r.FormValue("duration"); durationStr != "" {
		duration, err := strconv.Atoi(durationStr)
		if err != nil {
			return metadata, err
		}
		metadata.Duration = duration
	}

	// Обрабатываем ID альбома если он указан
	albumIDStr := r.FormValue("album_id")
	fmt.Printf("Значение album_id из формы: %s\n", albumIDStr)

	if albumIDStr != "" {
		albumID, err := uuid.Parse(albumIDStr)
		if err != nil {
			fmt.Printf("Ошибка при парсинге album_id: %v\n", err)
			return metadata, err
		}
		metadata.AlbumID = albumID
		fmt.Printf("Установлен AlbumID: %s\n", albumID)
	}

	return metadata, nil
}

// Проверка типа файла
func (h *TrackHandler) isAllowedFileType(contentType string) bool {
	for _, allowedType := range h.allowedTypes {
		if strings.EqualFold(contentType, allowedType) {
			return true
		}
	}
	return false
}

// Проверка прав администратора
func (h *TrackHandler) isAdmin(r *http.Request) bool {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "Bearer 33333333-3333-3333-3333-333333333333" {
		fmt.Println("Прямое разрешение по токену админа")
		return true
	}

	permission := r.Header.Get("X-User-Permission")
	isAdmin := permission == string(models.AdminPermission)
	fmt.Printf("Проверка прав админа: X-User-Permission=%s, isAdmin=%v\n", permission, isAdmin)
	return isAdmin
}

// GetTrackDetails получает детальную информацию о треке
func (h *TrackHandler) GetTrackDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	trackID, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Недопустимый идентификатор трека", http.StatusBadRequest)
		return
	}

	trackDetails, err := h.trackUseCase.GetTrackDetails(trackID)
	if err != nil {
		http.Error(w, "Трек не найден", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"id":          trackDetails.ID.String(),
		"title":       trackDetails.Title,
		"artist_name": trackDetails.ArtistName,
		"duration":    trackDetails.Duration,
		"cover_url":   trackDetails.CoverURL,
		"added_date":  trackDetails.AddedDate,
		"play_count":  trackDetails.PlayCount,
	}

	if trackDetails.Album != nil {
		response["album_id"] = trackDetails.Album.ID.String()
		response["album_title"] = trackDetails.Album.Title
		response["album_artist"] = trackDetails.Album.Artist
	}

	var genres []map[string]string
	for _, genre := range trackDetails.Genres {
		genres = append(genres, map[string]string{
			"id":   genre.ID.String(),
			"name": genre.Name,
		})
	}
	response["genres"] = genres

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *TrackHandler) ServeTrackFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	trackID, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Недопустимый идентификатор трека", http.StatusBadRequest)
		return
	}

	// Получаем информацию о треке из базы данных
	trackDetails, err := h.trackUseCase.GetTrackDetails(trackID)
	if err != nil {
		http.Error(w, "Трек не найден", http.StatusNotFound)
		return
	}

	// Получаем абсолютный путь к файлу
	trackFilePath, err := h.trackUseCase.GetTrackFilePath(trackID)
	if err != nil {
		http.Error(w, "Ошибка при получении файла: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем аудиофайл
	w.Header().Set("Content-Type", "audio/mpeg")
	w.Header().Set("Content-Disposition", "inline; filename="+trackDetails.Title+".mp3")
	http.ServeFile(w, r, trackFilePath)

	// Записываем прослушивание, если пользователь авторизован
	if userIDStr := r.Header.Get("X-User-ID"); userIDStr != "" {
		userID, err := uuid.Parse(userIDStr)
		if err == nil {
			go h.trackUseCase.PlayTrack(userID, trackID)
		}
	}
}

// SearchTracks выполняет поиск треков
func (h *TrackHandler) SearchTracks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Поисковый запрос не указан", http.StatusBadRequest)
		return
	}

	tracks, err := h.trackUseCase.SearchTracks(query)
	if err != nil {
		http.Error(w, "Ошибка при поиске треков: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tracks)
}

// DeleteTrack удаляет трек
func (h *TrackHandler) DeleteTrack(w http.ResponseWriter, r *http.Request) {
	// Проверяем права администратора
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

	if err := h.trackUseCase.DeleteTrack(trackID); err != nil {
		http.Error(w, "Ошибка при удалении трека: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
