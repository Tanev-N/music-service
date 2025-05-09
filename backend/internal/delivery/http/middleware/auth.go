package middleware

import (
	"context"
	"fmt"
	"log"
	"music-service/internal/usecases/interfaces"
	"net/http"
	"os"
	"strings"
)

// Список публичных маршрутов, не требующих авторизации
var publicRoutes = map[string]bool{
	"/api/v1/users":      true, // Регистрация
	"/api/v1/users/auth": true, // Аутентификация
	"/api/v1/tracks":     true, // Поиск треков (GET)
	"/api/v1/albums":     true, // Список альбомов (GET)
	"/api/v1/genres":     true, // Список жанров (GET)
}

// isPublicRoute проверяет, является ли маршрут публичным
func isPublicRoute(path string, method string) bool {
	fmt.Fprintf(os.Stderr, "Проверка маршрута: %s %s\n", method, path)

	if publicRoutes[path] {
		return true
	}

	if method == "GET" {
		if path == "/api/v1/tracks" {
			return true
		}

		if strings.Contains(path, "/stream") {
			return true
		}

		if strings.HasPrefix(path, "/api/v1/tracks/") {
			return true
		}

		if strings.HasPrefix(path, "/api/v1/albums") {
			return true
		}

		if strings.HasPrefix(path, "/api/v1/genres/tracks/") {
			return true
		}
	}

	if method == "OPTIONS" {
		return true
	}

	fmt.Fprintf(os.Stderr, "Маршрут не публичный: %s %s\n", method, path)
	return false
}

// AuthMiddleware извлекает информацию о пользователе из токена и добавляет ее в заголовки
func AuthMiddleware(userUseCase interfaces.UserUseCase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if isPublicRoute(r.URL.Path, r.Method) {
				next.ServeHTTP(w, r)
				return
			}

			log.Printf("Проверка авторизации для маршрута: %s %s", r.Method, r.URL.Path)
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Printf("Отсутствует заголовок Authorization")
				http.Error(w, "Необходима авторизация", http.StatusUnauthorized)
				return
			}

			log.Printf("Получен заголовок Authorization: %s", authHeader)
			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				log.Printf("Неверный формат токена: %v", tokenParts)
				http.Error(w, "Неверный формат токена авторизации", http.StatusUnauthorized)
				return
			}
			token := tokenParts[1]
			log.Printf("Проверка токена: %s", token)

			user, err := userUseCase.ValidateSession(token)
			if err != nil {
				log.Printf("Ошибка валидации токена: %v", err)
				http.Error(w, "Недействительный токен авторизации", http.StatusUnauthorized)
				return
			}

			log.Printf("Токен валиден, пользователь: %s (ID: %s)", user.Login, user.ID)
			r.Header.Set("X-User-ID", user.ID.String())
			r.Header.Set("X-User-Permission", string(user.Permission))

			ctx := context.WithValue(r.Context(), "userID", user.ID)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
