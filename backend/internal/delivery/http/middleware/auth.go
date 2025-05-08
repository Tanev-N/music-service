package middleware

import (
	"fmt"
	"music-service/internal/usecases/interfaces"
	"net/http"
	"os"
	"strings"
)

// AuthMiddleware извлекает информацию о пользователе из токена и добавляет ее в заголовки
func AuthMiddleware(userUseCase interfaces.UserUseCase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(os.Stderr, "AUTH MIDDLEWARE: запрос %s %s\n", r.Method, r.URL.Path)

			// Извлекаем токен из заголовка Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				fmt.Fprintf(os.Stderr, "AUTH MIDDLEWARE: заголовок Authorization отсутствует\n")
				next.ServeHTTP(w, r)
				return
			}

			// Извлекаем токен из заголовка
			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				fmt.Fprintf(os.Stderr, "AUTH MIDDLEWARE: неверный формат токена: %s\n", authHeader)
				next.ServeHTTP(w, r)
				return
			}
			token := tokenParts[1]
			fmt.Fprintf(os.Stderr, "AUTH MIDDLEWARE: получен токен: %s\n", token)

			// Пытаемся получить пользователя по токену
			user, err := userUseCase.ValidateSession(token)
			if err != nil {
				fmt.Fprintf(os.Stderr, "AUTH MIDDLEWARE: ошибка при проверке токена: %v\n", err)
				next.ServeHTTP(w, r)
				return
			}

			// Добавляем информацию о пользователе в заголовки
			r.Header.Set("X-User-ID", user.ID.String())
			r.Header.Set("X-User-Permission", string(user.Permission))

			fmt.Fprintf(os.Stderr, "AUTH MIDDLEWARE: аутентификация успешна: пользователь=%s, права=%s\n",
				user.Login, user.Permission)

			next.ServeHTTP(w, r)
		})
	}
}
