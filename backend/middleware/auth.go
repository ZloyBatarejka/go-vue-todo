package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"goTodo/backend/models"
	"goTodo/backend/services"
)

type contextKey string

const userIDContextKey contextKey = "userId"

// AuthMiddleware проверяет Bearer access-токен и пускает только авторизованные запросы.
// Middleware ожидает заголовок Authorization в формате "Bearer <token>".
// После валидации записывает userId из JWT в context запроса.
func AuthMiddleware(authService services.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
			if authHeader == "" {
				respondWithUnauthorized(w, "Authorization header is required")
				return
			}

			const bearerPrefix = "Bearer "
			if !strings.HasPrefix(authHeader, bearerPrefix) {
				respondWithUnauthorized(w, "Authorization header must be in format: Bearer <token>")
				return
			}

			token := strings.TrimSpace(strings.TrimPrefix(authHeader, bearerPrefix))
			if token == "" {
				respondWithUnauthorized(w, "Access token is required")
				return
			}

			claims, err := authService.ValidateAccessToken(token)
			if err != nil {
				respondWithUnauthorized(w, "Invalid or expired access token")
				return
			}
			if claims.UserID <= 0 {
				respondWithUnauthorized(w, "Invalid access token payload")
				return
			}

			ctx := context.WithValue(r.Context(), userIDContextKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// UserIDFromContext извлекает userId, который ранее положил AuthMiddleware.
// Возвращает userId и признак успешного извлечения.
// Если middleware не применен, ok будет false.
func UserIDFromContext(ctx context.Context) (int64, bool) {
	value := ctx.Value(userIDContextKey)
	userID, ok := value.(int64)
	if !ok || userID <= 0 {
		return 0, false
	}

	return userID, true
}

func respondWithUnauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)

	_ = json.NewEncoder(w).Encode(models.ErrorResponse{Error: message})
}
