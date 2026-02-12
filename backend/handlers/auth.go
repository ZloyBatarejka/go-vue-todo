package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/lib/pq"

	"goTodo/backend/models"
	"goTodo/backend/repository"
	"goTodo/backend/services"
)

const pgUniqueViolationCode = "23505"

// AuthHandler обрабатывает публичные auth-эндпоинты (register/login).
// Внутри использует UserRepository для доступа к пользователям и AuthService
// для хеширования пароля и работы с JWT.
type AuthHandler struct {
	userRepo repository.UserRepository
	auth     services.AuthService
}

// NewAuthHandler создает новый обработчик для auth-эндпоинтов.
// Параметры: userRepo — слой доступа к users; auth — сервис bcrypt/JWT.
// Возвращает: инициализированный AuthHandler.
func NewAuthHandler(userRepo repository.UserRepository, auth services.AuthService) *AuthHandler {
	return &AuthHandler{
		userRepo: userRepo,
		auth:     auth,
	}
}

// Register godoc
// @Summary Register user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Register request"
// @Success 201 {object} models.AuthResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 409 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	username := strings.TrimSpace(req.Username)
	if username == "" || req.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Fields 'username' and 'password' are required")
		return
	}

	passwordHash, err := h.auth.HashPassword(req.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to register user")
		return
	}

	user, err := h.userRepo.CreateUser(username, passwordHash)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolationCode {
			respondWithError(w, http.StatusConflict, "Username is already taken")
			return
		}

		log.Printf("Error creating user: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to register user")
		return
	}

	accessToken, err := h.auth.GenerateAccessToken(user.ID, user.Username)
	if err != nil {
		log.Printf("Error generating access token: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to register user")
		return
	}

	respondWithJSON(w, http.StatusCreated, models.AuthResponse{
		AccessToken: accessToken,
		User:        toUserResponse(user),
	})
}

// Login godoc
// @Summary Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login request"
// @Success 200 {object} models.AuthResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	username := strings.TrimSpace(req.Username)
	if username == "" || req.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Fields 'username' and 'password' are required")
		return
	}

	user, err := h.userRepo.FindByUsername(username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusUnauthorized, "Invalid username or password")
			return
		}

		log.Printf("Error finding user: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to login")
		return
	}

	if err := h.auth.VerifyPassword(req.Password, user.PasswordHash); err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			respondWithError(w, http.StatusUnauthorized, "Invalid username or password")
			return
		}

		log.Printf("Error verifying password: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to login")
		return
	}

	accessToken, err := h.auth.GenerateAccessToken(user.ID, user.Username)
	if err != nil {
		log.Printf("Error generating access token: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to login")
		return
	}

	respondWithJSON(w, http.StatusOK, models.AuthResponse{
		AccessToken: accessToken,
		User:        toUserResponse(user),
	})
}

func toUserResponse(user *models.User) models.UserResponse {
	return models.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}
}
