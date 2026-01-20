// @title goTodo API
// @version 1.0
// @host localhost:8080
// @BasePath /api
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"

	"goTodo/backend/database"
	"goTodo/backend/handlers"
	"goTodo/backend/repository"

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

	todoHandler := handlers.NewTodoHandler(todoRepo)

	router := mux.NewRouter()

	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/todos", todoHandler.GetAllTodos).Methods("GET")
	api.HandleFunc("/todos", todoHandler.CreateTodo).Methods("POST")
	api.HandleFunc("/todos/{id:[0-9]+}", todoHandler.GetTodo).Methods("GET")
	api.HandleFunc("/todos/{id:[0-9]+}", todoHandler.DeleteTodo).Methods("DELETE")

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	port := ":8080"
	fmt.Printf("Server starting on http://localhost%s\n", port)
	fmt.Println("API endpoints:")
	fmt.Println("  GET    /api/todos")
	fmt.Println("  POST   /api/todos")
	fmt.Println("  GET    /api/todos/{id}")
	fmt.Println("  DELETE /api/todos/{id}")
	fmt.Println("Swagger UI:")
	fmt.Println("  GET    /swagger/index.html")

	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
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
