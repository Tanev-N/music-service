package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"music-service/internal/config"
	"music-service/internal/repository/db"
	"music-service/internal/repository/postgres"
	"music-service/internal/usecases"
	"music-service/src/ui"
)

func main() {
	cfg, err := config.NewConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("Ошибка при загрузке конфигурации: %v", err)
	}

	log.Printf("Конфигурация: host=%s, port=%s, user=%s, dbname=%s",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.Username, cfg.DB.DBName)

	db, err := db.NewPostgresDB(cfg.DB)
	if err != nil {
		log.Fatalf("Ошибка при подключении к базе данных: %v", err)
	}
	defer db.Close()

	userRepo := postgres.NewUserRepository(db)
	sessionRepo := postgres.NewSessionRepository(db)
	trackRepo := postgres.NewTrackRepository(db)
	albumRepo := postgres.NewAlbumRepository(db)
	playlistRepo := postgres.NewPlaylistRepository(db)
	genreRepo := postgres.NewGenreRepository(db)
	historyRepo := postgres.NewHistoryRepository(db)

	userUseCase := usecases.NewUserUseCase(userRepo, sessionRepo)
	trackUseCase := usecases.NewTrackUseCase(trackRepo, historyRepo, albumRepo)
	albumUseCase := usecases.NewAlbumUseCase(albumRepo, trackRepo)
	playlistUseCase := usecases.NewPlaylistUseCase(playlistRepo, trackRepo, userRepo)
	genreUseCase := usecases.NewGenreUseCase(genreRepo, trackRepo)
	historyUseCase := usecases.NewHistoryUseCase(historyRepo, trackRepo)

	app := ui.NewConsoleApp(
		userUseCase,
		trackUseCase,
		albumUseCase,
		playlistUseCase,
		genreUseCase,
		historyUseCase,
	)

	args := os.Args
	if len(args) > 1 {
		command := strings.Join(args[1:], " ")
		output, err := app.ExecuteCommand(command)
		if err != nil {
			fmt.Printf("Ошибка: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(output)
	} else {
		app.StartInteractive()
	}
}
