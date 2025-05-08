package main

import (
	"database/sql"
	"fmt"
	"log"
	"music-service/internal/config"
	"music-service/internal/delivery/http/router"
	"music-service/internal/repository/postgres"
	"music-service/internal/usecases"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.NewConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	userRepo := postgres.NewUserRepository(db)
	sessionRepo := postgres.NewSessionRepository(db)
	userUseCase := usecases.NewUserUseCase(userRepo, sessionRepo)

	r := router.NewRouter(userUseCase)

	port := ":" + cfg.App.Port
	fmt.Printf("Сервер запущен на http://localhost%s\n", port)

	if err := http.ListenAndServe(port, r.GetRouter()); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
