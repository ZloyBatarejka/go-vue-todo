package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidToken           = errors.New("invalid access token")
	ErrInvalidTokenSigning    = errors.New("invalid token signing method")
	ErrInvalidCredentials     = errors.New("invalid credentials")
	ErrEmptyJWTSecret         = errors.New("jwt secret is required")
	ErrInvalidAccessTokenTTL  = errors.New("access token ttl must be greater than zero")
	ErrInvalidRefreshTokenTTL = errors.New("refresh token ttl must be greater than zero")
	ErrFailedToGenerateToken  = errors.New("failed to generate token")
)

// AuthService описывает операции прикладной авторизации.
// Сервис инкапсулирует работу с bcrypt и JWT, чтобы хендлеры не знали
// деталей хеширования, подписи токенов и валидации claims.
type AuthService interface {
	HashPassword(password string) (string, error)
	VerifyPassword(password string, passwordHash string) error
	GenerateAccessToken(userID int64, username string) (string, error)
	ValidateAccessToken(token string) (*AccessTokenClaims, error)
	GenerateRefreshToken() (string, string, error)
	HashRefreshToken(token string) string
}

// AccessTokenClaims хранит payload access-токена.
// Включает пользовательские поля (UserID, Username) и стандартные
// зарегистрированные JWT claims (sub, iat, exp).
type AccessTokenClaims struct {
	UserID   int64  `json:"userId"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// authService — конкретная реализация AuthService.
// Содержит секрет подписи, TTL access-токена и источник времени,
// который можно подменять в тестах для предсказуемых сценариев.
type authService struct {
	jwtSecret       []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	now             func() time.Time
}

// NewAuthService создает новый экземпляр AuthService с проверкой входных параметров.
// Параметры: jwtSecret — секрет для HS256 подписи; accessTokenTTL — срок жизни access.
// Возвращает: готовый сервис или ошибку, если секрет пустой/TTL некорректный.
func NewAuthService(jwtSecret string, accessTokenTTL time.Duration, refreshTokenTTL time.Duration) (AuthService, error) {
	if strings.TrimSpace(jwtSecret) == "" {
		return nil, ErrEmptyJWTSecret
	}
	if accessTokenTTL <= 0 {
		return nil, ErrInvalidAccessTokenTTL
	}
	if refreshTokenTTL <= 0 {
		return nil, ErrInvalidRefreshTokenTTL
	}

	return &authService{
		jwtSecret:       []byte(jwtSecret),
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
		now:             time.Now,
	}, nil
}

// HashPassword хеширует открытый пароль через bcrypt.
// Параметры: password — пароль в открытом виде.
// Возвращает: строку хеша для БД или ошибку, если хеширование не удалось.
func (s *authService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}

// VerifyPassword сверяет пароль пользователя с bcrypt-хешем из БД.
// Параметры: password — введенный пароль; passwordHash — сохраненный хеш.
// Возвращает: nil при совпадении, ErrInvalidCredentials при mismatch, либо тех. ошибку.
func (s *authService) VerifyPassword(password string, passwordHash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrInvalidCredentials
		}
		return fmt.Errorf("failed to verify password: %w", err)
	}
	return nil
}

// GenerateAccessToken выпускает подписанный access JWT с пользовательскими claims.
// Параметры: userID — идентификатор пользователя; username — логин пользователя.
// Возвращает: строку JWT или ошибку, если токен не удалось подписать.
func (s *authService) GenerateAccessToken(userID int64, username string) (string, error) {
	now := s.now().UTC()
	claims := AccessTokenClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(userID, 10),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTokenTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign access token: %w", err)
	}

	return tokenString, nil
}

// ValidateAccessToken валидирует access JWT и извлекает claims.
// Параметры: token — строка JWT из заголовка Authorization (без префикса Bearer).
// Возвращает: AccessTokenClaims при успехе или ErrInvalidToken/другую ошибку.
func (s *authService) ValidateAccessToken(token string) (*AccessTokenClaims, error) {
	claims := &AccessTokenClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidTokenSigning
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}
	if !parsedToken.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// GenerateRefreshToken создает новый opaque refresh token и его SHA-256 hash для хранения в БД.
func (s *authService) GenerateRefreshToken() (string, string, error) {
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", "", fmt.Errorf("%w: %v", ErrFailedToGenerateToken, err)
	}

	rawToken := base64.RawURLEncoding.EncodeToString(randomBytes)
	return rawToken, s.HashRefreshToken(rawToken), nil
}

func (s *authService) HashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
