package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"music-service/internal/models"
	"music-service/internal/usecases/interfaces"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type PlaylistHandler struct {
	playlistUseCase interfaces.PlaylistUseCase
	userUseCase     interfaces.UserUseCase
}

func NewPlaylistHandler(
	playlistUseCase interfaces.PlaylistUseCase,
	userUseCase interfaces.UserUseCase,
) *PlaylistHandler {
	return &PlaylistHandler{
		playlistUseCase: playlistUseCase,
		userUseCase:     userUseCase,
	}
}

// CreatePlaylist создает новый плейлист
func (h *PlaylistHandler) CreatePlaylist(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		CoverURL    string `json:"cover_url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	userID, err := getUserIDFromSession(r)
	if err != nil {
		http.Error(w, "Необходима авторизация", http.StatusUnauthorized)
		return
	}

	playlist, err := h.playlistUseCase.CreatePlaylist(userID, request.Name, request.Description, request.CoverURL)
	if err != nil {
		log.Printf("Ошибка при создании плейлиста: %v", err)
		http.Error(w, "Ошибка при создании плейлиста", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(playlist)
}

// GetPlaylistWithTracks возвращает плейлист с треками
func (h *PlaylistHandler) GetPlaylistWithTracks(w http.ResponseWriter, r *http.Request) {
	playlistID, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Неверный ID плейлиста", http.StatusBadRequest)
		return
	}

	playlistWithTracks, err := h.playlistUseCase.GetPlaylistWithTracks(playlistID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "Плейлист не найден", http.StatusNotFound)
		} else {
			log.Printf("Ошибка при получении плейлиста: %v", err)
			http.Error(w, "Ошибка при получении плейлиста", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(playlistWithTracks)
}

// EditPlaylistInfo обновляет информацию о плейлисте
func (h *PlaylistHandler) EditPlaylistInfo(w http.ResponseWriter, r *http.Request) {
	playlistID, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Неверный ID плейлиста", http.StatusBadRequest)
		return
	}

	var request struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	// Проверяем, что пользователь является владельцем плейлиста
	userID, err := getUserIDFromSession(r)
	if err != nil {
		http.Error(w, "Необходима авторизация", http.StatusUnauthorized)
		return
	}

	playlist, err := h.playlistUseCase.GetPlaylistWithTracks(playlistID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "Плейлист не найден", http.StatusNotFound)
		} else {
			log.Printf("Ошибка при получении плейлиста: %v", err)
			http.Error(w, "Ошибка при получении плейлиста", http.StatusInternalServerError)
		}
		return
	}

	// Проверяем, что пользователь является владельцем плейлиста
	if playlist.Playlist.UserID != userID {
		http.Error(w, "Нет доступа к редактированию плейлиста", http.StatusForbidden)
		return
	}

	if err := h.playlistUseCase.EditPlaylistInfo(playlistID, request.Name, request.Description); err != nil {
		log.Printf("Ошибка при обновлении плейлиста: %v", err)
		http.Error(w, "Ошибка при обновлении плейлиста", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Информация о плейлисте обновлена"}`))
}

// GetPlaylistTracks возвращает треки из плейлиста
func (h *PlaylistHandler) GetPlaylistTracks(w http.ResponseWriter, r *http.Request) {
	playlistID, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Неверный ID плейлиста", http.StatusBadRequest)
		return
	}

	tracks, err := h.playlistUseCase.GetPlaylistTracks(playlistID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "Плейлист не найден", http.StatusNotFound)
		} else {
			log.Printf("Ошибка при получении треков: %v", err)
			http.Error(w, "Ошибка при получении треков", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tracks)
}

// AddTrackToPlaylist добавляет трек в плейлист
func (h *PlaylistHandler) AddTrackToPlaylist(w http.ResponseWriter, r *http.Request) {
	playlistID, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Неверный ID плейлиста", http.StatusBadRequest)
		return
	}

	var request struct {
		TrackID uuid.UUID `json:"track_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	userID, err := getUserIDFromSession(r)
	if err != nil {
		http.Error(w, "Необходима авторизация", http.StatusUnauthorized)
		return
	}

	playlist, err := h.playlistUseCase.GetPlaylistWithTracks(playlistID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "Плейлист не найден", http.StatusNotFound)
		} else {
			log.Printf("Ошибка при получении плейлиста: %v", err)
			http.Error(w, "Ошибка при получении плейлиста", http.StatusInternalServerError)
		}
		return
	}

	if playlist.Playlist.UserID != userID {
		http.Error(w, "Нет доступа к редактированию плейлиста", http.StatusForbidden)
		return
	}

	if err := h.playlistUseCase.AddTrackToPlaylist(playlistID, request.TrackID); err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "Трек не найден", http.StatusNotFound)
		} else {
			log.Printf("Ошибка при добавлении трека: %v", err)
			http.Error(w, "Ошибка при добавлении трека", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Трек добавлен в плейлист"}`))
}

// RemoveTrackFromPlaylist удаляет трек из плейлиста
func (h *PlaylistHandler) RemoveTrackFromPlaylist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	playlistID, err := uuid.Parse(vars["playlistId"])
	if err != nil {
		http.Error(w, "Неверный ID плейлиста", http.StatusBadRequest)
		return
	}

	trackID, err := uuid.Parse(vars["trackId"])
	if err != nil {
		http.Error(w, "Неверный ID трека", http.StatusBadRequest)
		return
	}

	userID, err := getUserIDFromSession(r)
	if err != nil {
		http.Error(w, "Необходима авторизация", http.StatusUnauthorized)
		return
	}

	playlist, err := h.playlistUseCase.GetPlaylistWithTracks(playlistID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "Плейлист не найден", http.StatusNotFound)
		} else {
			log.Printf("Ошибка при получении плейлиста: %v", err)
			http.Error(w, "Ошибка при получении плейлиста", http.StatusInternalServerError)
		}
		return
	}

	if playlist.Playlist.UserID != userID {
		http.Error(w, "Нет доступа к редактированию плейлиста", http.StatusForbidden)
		return
	}

	if err := h.playlistUseCase.RemoveTrackFromPlaylist(playlistID, trackID); err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "Трек не найден в плейлисте", http.StatusNotFound)
		} else {
			log.Printf("Ошибка при удалении трека: %v", err)
			http.Error(w, "Ошибка при удалении трека", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Трек удален из плейлиста"}`))
}

// GetUserPlaylists возвращает список плейлистов пользователя
func (h *PlaylistHandler) GetUserPlaylists(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromSession(r)
	if err != nil {
		http.Error(w, "Необходима авторизация", http.StatusUnauthorized)
		return
	}

	playlists, err := h.playlistUseCase.GetUserPlaylists(userID)
	if err != nil {
		log.Printf("Ошибка при получении плейлистов пользователя: %v", err)
		http.Error(w, "Ошибка при получении плейлистов", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(playlists)
}

// DeletePlaylist удаляет плейлист
func (h *PlaylistHandler) DeletePlaylist(w http.ResponseWriter, r *http.Request) {
	playlistID, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Неверный ID плейлиста", http.StatusBadRequest)
		return
	}

	userID, err := getUserIDFromSession(r)
	if err != nil {
		http.Error(w, "Необходима авторизация", http.StatusUnauthorized)
		return
	}

	if err := h.playlistUseCase.DeletePlaylist(playlistID, userID); err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "Плейлист не найден", http.StatusNotFound)
		} else if err.Error() == "user is not the owner of the playlist" {
			http.Error(w, "Нет доступа к удалению плейлиста", http.StatusForbidden)
		} else {
			log.Printf("Ошибка при удалении плейлиста: %v", err)
			http.Error(w, "Ошибка при удалении плейлиста", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Плейлист успешно удален"}`))
}

// getUserIDFromSession извлекает ID пользователя из сессии
func getUserIDFromSession(r *http.Request) (uuid.UUID, error) {
	fmt.Printf("Попытка получить ID пользователя из контекста запроса\n")

	// Проверяем наличие заголовка X-User-ID (для совместимости)
	userIDHeader := r.Header.Get("X-User-ID")
	if userIDHeader != "" {
		fmt.Printf("Найден заголовок X-User-ID: %s\n", userIDHeader)
		userID, err := uuid.Parse(userIDHeader)
		if err == nil {
			return userID, nil
		}
		fmt.Printf("Ошибка при парсинге ID из заголовка: %v\n", err)
	}

	// Пытаемся получить ID из контекста
	userID, ok := r.Context().Value("userID").(uuid.UUID)
	if !ok {
		fmt.Printf("ID пользователя не найден в контексте запроса\n")
		return uuid.Nil, errors.New("user ID not found in session")
	}

	fmt.Printf("Успешно получен ID пользователя из контекста: %s\n", userID.String())
	return userID, nil
}
