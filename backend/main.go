// @title goTodo API
// @version 1.0
// @host localhost:8080
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"

	"goTodo/backend/database"
	"goTodo/backend/handlers"
	"goTodo/backend/middleware"
	"goTodo/backend/repository"
	"goTodo/backend/services"

	_ "goTodo/backend/docs"
)

func main() {
	// Загружаем переменные окружения из .env (если файл есть)
	_ = godotenv.Load()

	db, err := database.NewDB(database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnvInt("DB_PORT", 5432),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", ""),
		DBName:   getEnv("DB_NAME", "postgres"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	todoRepo := repository.NewTodoRepository(db)
	userRepo := repository.NewUserRepository(db)
	refreshSessionRepo := repository.NewRefreshSessionRepository(db)

	refreshTokenTTL := time.Duration(getEnvInt("JWT_REFRESH_TTL_HOURS", 168)) * time.Hour

	authService, err := services.NewAuthService(
		getEnv("JWT_SECRET", "dev-secret-change-me"),
		time.Duration(getEnvInt("JWT_ACCESS_TTL_MINUTES", 60))*time.Minute,
		refreshTokenTTL,
	)
	if err != nil {
		log.Fatalf("Failed to initialize auth service: %v", err)
	}

	todoHandler := handlers.NewTodoHandler(todoRepo)
	authHandler := handlers.NewAuthHandler(
		userRepo,
		refreshSessionRepo,
		authService,
		refreshTokenTTL,
		handlers.RefreshCookieConfig{
			Name:     getEnv("REFRESH_COOKIE_NAME", "goTodo_refresh_token"),
			Domain:   getEnv("REFRESH_COOKIE_DOMAIN", ""),
			Path:     getEnv("REFRESH_COOKIE_PATH", "/api/auth"),
			Secure:   getEnvBool("REFRESH_COOKIE_SECURE", false),
			HTTPOnly: getEnvBool("REFRESH_COOKIE_HTTPONLY", true),
			SameSite: parseSameSite(getEnv("REFRESH_COOKIE_SAMESITE", "Lax")),
		},
	)

	router := mux.NewRouter()

	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/auth/register", authHandler.Register).Methods("POST")
	api.HandleFunc("/auth/login", authHandler.Login).Methods("POST")
	api.HandleFunc("/auth/refresh", authHandler.Refresh).Methods("POST")
	api.HandleFunc("/auth/logout", authHandler.Logout).Methods("POST")

	// Все todo-эндпоинты требуют валидный Bearer access-токен.
	authRequired := middleware.AuthMiddleware(authService)
	api.Handle("/todos", authRequired(http.HandlerFunc(todoHandler.GetAllTodos))).Methods("GET")
	api.Handle("/todos", authRequired(http.HandlerFunc(todoHandler.CreateTodo))).Methods("POST")
	api.Handle("/todos/{id:[0-9]+}", authRequired(http.HandlerFunc(todoHandler.GetTodo))).Methods("GET")
	api.Handle("/todos/{id:[0-9]+}", authRequired(http.HandlerFunc(todoHandler.DeleteTodo))).Methods("DELETE")

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	port := ":8080"
	fmt.Printf("Server starting on http://localhost%s\n", port)
	fmt.Println("API endpoints:")
	fmt.Println("  POST   /api/auth/register")
	fmt.Println("  POST   /api/auth/login")
	fmt.Println("  POST   /api/auth/refresh")
	fmt.Println("  POST   /api/auth/logout")
	fmt.Println("  GET    /api/todos")
	fmt.Println("  POST   /api/todos")
	fmt.Println("  GET    /api/todos/{id}")
	fmt.Println("  DELETE /api/todos/{id}")
	fmt.Println("Swagger UI:")
	fmt.Println("  GET    /swagger/index.html")

	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{getEnv("CORS_ALLOWED_ORIGIN", "http://localhost:5173")},
		AllowedMethods: []string{"GET", "POST", "DELETE", "OPTIONS"},
		// Authorization нужен для Bearer JWT; Cookie/Set-Cookie — для refresh flow.
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		ExposedHeaders:   []string{"Set-Cookie"},
		AllowCredentials: true,
	}).Handler(router)

	if err := http.ListenAndServe(port, corsHandler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func getEnv(key string, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return i
}

func getEnvBool(key string, fallback bool) bool {
	value := strings.TrimSpace(strings.ToLower(os.Getenv(key)))
	if value == "" {
		return fallback
	}

	switch value {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return fallback
	}
}

func parseSameSite(value string) http.SameSite {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "none":
		return http.SameSiteNoneMode
	case "strict":
		return http.SameSiteStrictMode
	case "lax":
		return http.SameSiteLaxMode
	default:
		return http.SameSiteLaxMode
	}
}
