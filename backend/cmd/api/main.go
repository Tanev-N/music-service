package main

import (
	"fmt"
	"log"
	"music-service/internal/config"
	"music-service/internal/delivery/http/router"
	"music-service/internal/repository"
	"music-service/internal/repository/db"
	"music-service/internal/usecases"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.NewConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	if err := os.MkdirAll(cfg.Storage.TracksDir, 0755); err != nil {
		log.Fatalf("Ошибка создания директории для хранения треков: %v", err)
	}

	tracksDir, err := filepath.Abs(cfg.Storage.TracksDir)
	if err != nil {
		log.Fatalf("Ошибка получения абсолютного пути к директории треков: %v", err)
	}

	dbConfig := db.Config{
		Host:      os.Getenv("DB_HOST"),
		Port:      os.Getenv("DB_PORT"),
		Username:  os.Getenv("DB_USER"),
		Password:  os.Getenv("DB_PASSWORD"),
		DBName:    os.Getenv("DB_NAME"),
		SSLMode:   os.Getenv("DB_SSLMODE"),
		TracksDir: tracksDir,
	}

	repo, err := repository.NewRepository(dbConfig)
	if err != nil {
		log.Fatalf("Ошибка создания репозитория: %v", err)
	}

	userUseCase := usecases.NewUserUseCase(repo.User, repo.Session)
	trackUseCase := usecases.NewTrackUseCase(
		repo.Track,
		repo.History,
		repo.Album,
		cfg.Storage.MaxFileSizeMB,
		cfg.Storage.AllowedTypes,
	)
	albumUseCase := usecases.NewAlbumUseCase(
		repo.Album,
		repo.Track,
	)
	genreUseCase := usecases.NewGenreUseCase(
		repo.Genre,
		repo.Track,
	)

	r := router.NewRouter(
		userUseCase,
		trackUseCase,
		albumUseCase,
		genreUseCase,
		cfg.Storage.AllowedTypes,
		cfg.Storage.MaxFileSizeMB,
	)

	port := ":" + cfg.App.Port
	fmt.Printf("Сервер запущен на http://localhost%s\n", port)

	if err := http.ListenAndServe(port, r.GetRouter()); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
