package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
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
	userRepo      repository.UserRepository
	refreshRepo   repository.RefreshSessionRepository
	auth          services.AuthService
	refreshTTL    time.Duration
	refreshCookie RefreshCookieConfig
}

type RefreshCookieConfig struct {
	Name     string
	Domain   string
	Path     string
	Secure   bool
	HTTPOnly bool
	SameSite http.SameSite
}

// NewAuthHandler создает новый обработчик для auth-эндпоинтов.
// Параметры: userRepo — слой доступа к users; auth — сервис bcrypt/JWT.
// Возвращает: инициализированный AuthHandler.
func NewAuthHandler(
	userRepo repository.UserRepository,
	refreshRepo repository.RefreshSessionRepository,
	auth services.AuthService,
	refreshTTL time.Duration,
	refreshCookie RefreshCookieConfig,
) *AuthHandler {
	return &AuthHandler{
		userRepo:      userRepo,
		refreshRepo:   refreshRepo,
		auth:          auth,
		refreshTTL:    refreshTTL,
		refreshCookie: refreshCookie,
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

	if err := h.issueRefreshSessionAndSetCookie(w, user.ID); err != nil {
		log.Printf("Error issuing refresh session on register: %v", err)
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

	if err := h.issueRefreshSessionAndSetCookie(w, user.ID); err != nil {
		log.Printf("Error issuing refresh session on login: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to login")
		return
	}

	respondWithJSON(w, http.StatusOK, models.AuthResponse{
		AccessToken: accessToken,
		User:        toUserResponse(user),
	})
}

// Refresh godoc
// @Summary Refresh access token
// @Tags auth
// @Produce json
// @Success 200 {object} models.AuthResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := h.readRefreshCookie(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Refresh token is required")
		return
	}

	hashedToken := h.auth.HashRefreshToken(refreshToken)
	session, err := h.refreshRepo.FindByTokenHash(hashedToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.clearRefreshCookie(w)
			respondWithError(w, http.StatusUnauthorized, "Invalid refresh token")
			return
		}

		log.Printf("Error reading refresh session: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to refresh session")
		return
	}

	if session.RevokedAt != nil || session.ConsumedAt != nil || session.ReplacedBySessionID != nil || session.ExpiresAt.Before(time.Now().UTC()) {
		_ = h.refreshRepo.RevokeFamily(session.FamilyID, "refresh token reuse or expired token")
		h.clearRefreshCookie(w)
		respondWithError(w, http.StatusUnauthorized, "Refresh token is not active")
		return
	}

	user, err := h.userRepo.FindByID(session.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.clearRefreshCookie(w)
			respondWithError(w, http.StatusUnauthorized, "User not found")
			return
		}

		log.Printf("Error finding user during refresh: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to refresh session")
		return
	}

	accessToken, err := h.auth.GenerateAccessToken(user.ID, user.Username)
	if err != nil {
		log.Printf("Error generating access token on refresh: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to refresh session")
		return
	}

	newRefreshToken, newRefreshHash, err := h.auth.GenerateRefreshToken()
	if err != nil {
		log.Printf("Error generating refresh token on refresh: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to refresh session")
		return
	}

	if err := h.refreshRepo.RotateSession(
		session.ID,
		session.UserID,
		session.FamilyID,
		newRefreshHash,
		time.Now().UTC().Add(h.refreshTTL),
	); err != nil {
		log.Printf("Error rotating refresh token: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to refresh session")
		return
	}

	h.setRefreshCookie(w, newRefreshToken, time.Now().UTC().Add(h.refreshTTL))
	respondWithJSON(w, http.StatusOK, models.AuthResponse{
		AccessToken: accessToken,
		User:        toUserResponse(user),
	})
}

// Logout godoc
// @Summary Logout user
// @Tags auth
// @Success 204 "No Content"
// @Failure 500 {object} models.ErrorResponse
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := h.readRefreshCookie(r)
	if err == nil && strings.TrimSpace(refreshToken) != "" {
		tokenHash := h.auth.HashRefreshToken(refreshToken)
		if revokeErr := h.refreshRepo.RevokeByTokenHash(tokenHash, "user logout"); revokeErr != nil {
			log.Printf("Error revoking refresh session on logout: %v", revokeErr)
			respondWithError(w, http.StatusInternalServerError, "Failed to logout")
			return
		}
	}

	h.clearRefreshCookie(w)
	w.WriteHeader(http.StatusNoContent)
}

func toUserResponse(user *models.User) models.UserResponse {
	return models.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}
}

func (h *AuthHandler) issueRefreshSessionAndSetCookie(w http.ResponseWriter, userID int64) error {
	refreshToken, refreshHash, err := h.auth.GenerateRefreshToken()
	if err != nil {
		return err
	}

	familyID := uuid.NewString()
	expiresAt := time.Now().UTC().Add(h.refreshTTL)
	if _, err := h.refreshRepo.CreateSession(userID, familyID, refreshHash, expiresAt); err != nil {
		return err
	}

	h.setRefreshCookie(w, refreshToken, expiresAt)
	return nil
}

func (h *AuthHandler) setRefreshCookie(w http.ResponseWriter, token string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     h.refreshCookie.Name,
		Value:    token,
		Path:     h.refreshCookie.Path,
		Domain:   h.refreshCookie.Domain,
		Expires:  expiresAt,
		HttpOnly: h.refreshCookie.HTTPOnly,
		Secure:   h.refreshCookie.Secure,
		SameSite: h.refreshCookie.SameSite,
	})
}

func (h *AuthHandler) clearRefreshCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     h.refreshCookie.Name,
		Value:    "",
		Path:     h.refreshCookie.Path,
		Domain:   h.refreshCookie.Domain,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: h.refreshCookie.HTTPOnly,
		Secure:   h.refreshCookie.Secure,
		SameSite: h.refreshCookie.SameSite,
	})
}

func (h *AuthHandler) readRefreshCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie(h.refreshCookie.Name)
	if err != nil {
		return "", err
	}

	token := strings.TrimSpace(cookie.Value)
	if token == "" {
		return "", http.ErrNoCookie
	}

	return token, nil
}
