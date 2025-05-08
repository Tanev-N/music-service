package usecases

import (
	"errors"
	"fmt"
	"music-service/internal/models"
	"music-service/internal/repository/interfaces"
	usecaseInterfaces "music-service/internal/usecases/interfaces"
	"sort"
	"time"

	"github.com/google/uuid"
)

type historyUseCase struct {
	historyRepo interfaces.HistoryRepository
	trackRepo   interfaces.TrackRepository
}

func NewHistoryUseCase(
	historyRepo interfaces.HistoryRepository,
	trackRepo interfaces.TrackRepository,
) usecaseInterfaces.HistoryUseCase {
	return &historyUseCase{
		historyRepo: historyRepo,
		trackRepo:   trackRepo,
	}
}

func (uc *historyUseCase) RecordPlayback(userID uuid.UUID, trackID uuid.UUID) error {
	track, err := uc.trackRepo.FindByID(trackID)
	if err != nil {
		return fmt.Errorf("track not found: %w", err)
	}
	if track.Duration < 30 {
		return errors.New("track is too short to record playback")
	}

	history, err := uc.historyRepo.GetHistory(userID)
	if err != nil {
		return fmt.Errorf("failed to get user history: %w", err)
	}

	recentPlays := 0
	for _, entry := range history {
		if entry.TrackID == trackID && time.Since(entry.ListenedAt) < 5*time.Minute {
			recentPlays++
			if recentPlays >= 3 {
				return errors.New("track played too frequently")
			}
		}
	}

	if err := uc.historyRepo.AddEntry(userID, trackID); err != nil {
		return fmt.Errorf("failed to record playback: %w", err)
	}

	return nil
}

func (uc *historyUseCase) GetUserHistory(userID uuid.UUID) ([]*models.ListeningHistory, error) {
	history, err := uc.historyRepo.GetHistory(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get listening history: %w", err)
	}

	sort.Slice(history, func(i, j int) bool {
		return history[i].ListenedAt.After(history[j].ListenedAt)
	})

	const maxHistoryItems = 100
	if len(history) > maxHistoryItems {
		history = history[:maxHistoryItems]
	}

	return history, nil
}

func (uc *historyUseCase) GetRecentPlays(userID uuid.UUID, within time.Duration) ([]*models.ListeningHistory, error) {
	history, err := uc.historyRepo.GetHistory(userID)
	if err != nil {
		return nil, err
	}

	var recent []*models.ListeningHistory
	now := time.Now()

	for _, entry := range history {
		if now.Sub(entry.ListenedAt) <= within {
			recent = append(recent, entry)
		}
	}

	return recent, nil
}
