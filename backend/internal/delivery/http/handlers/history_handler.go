package handlers

import (
	"encoding/json"
	"fmt"
	"music-service/internal/usecases/interfaces"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type HistoryHandler struct {
	historyUseCase interfaces.HistoryUseCase
}

func NewHistoryHandler(historyUseCase interfaces.HistoryUseCase) *HistoryHandler {
	return &HistoryHandler{
		historyUseCase: historyUseCase,
	}
}

// RecordPlayback записывает прослушивание трека
func (h *HistoryHandler) RecordPlayback(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	trackID, err := uuid.Parse(vars["trackId"])
	if err != nil {
		http.Error(w, "Некорректный ID трека", http.StatusBadRequest)
		return
	}

	userIDStr := r.Header.Get("X-User-ID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Не удалось определить пользователя", http.StatusUnauthorized)
		return
	}

	err = h.historyUseCase.RecordPlayback(userID, trackID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetUserHistory возвращает историю прослушиваний пользователя
func (h *HistoryHandler) GetUserHistory(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromSession(r)
	if err != nil {
		http.Error(w, "Не удалось определить пользователя", http.StatusUnauthorized)
		return
	}

	history, err := h.historyUseCase.GetUserHistory(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}

// GetRecentPlays возвращает недавние прослушивания
func (h *HistoryHandler) GetRecentPlays(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Header.Get("X-User-ID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Не удалось определить пользователя", http.StatusUnauthorized)
		return
	}

	hours := 24
	if hoursParam := r.URL.Query().Get("hours"); hoursParam != "" {
		fmt.Sscanf(hoursParam, "%d", &hours)
	}

	history, err := h.historyUseCase.GetRecentPlays(userID, time.Duration(hours)*time.Hour)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}
