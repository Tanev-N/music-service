package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type Config struct {
	Host      string
	Port      string
	Username  string
	Password  string
	DBName    string
	SSLMode   string
	TracksDir string
}

func NewPostgresDB(cfg Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("host='%s' port='%s' user='%s' password='%s' dbname='%s' sslmode='%s'",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)

	log.Printf("DSN строка подключения: %s", dsn)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Попытки подключения к БД с повторами
	maxRetries := 10
	retryDelay := 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		err = db.Ping()
		if err == nil {
			log.Println("Successfully connected to PostgreSQL database")
			return db, nil
		}

		log.Printf("Не удалось подключиться к БД (попытка %d/%d): %v", i+1, maxRetries, err)

		if i < maxRetries-1 {
			log.Printf("Ожидание %v перед следующей попыткой...", retryDelay)
			time.Sleep(retryDelay)
		}
	}

	return nil, fmt.Errorf("не удалось подключиться к БД после %d попыток: %w", maxRetries, err)
}
