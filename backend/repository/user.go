package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"goTodo/backend/models"
)

type UserRepository interface {
	CreateUser(username string, passwordHash string) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	FindByID(id int64) (*models.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(username string, passwordHash string) (*models.User, error) {
	user := &models.User{}
	query := `
		INSERT INTO users (username, password_hash)
		VALUES ($1, $2)
		RETURNING id, username, password_hash, created_at::text
	`

	err := r.db.QueryRow(query, username, passwordHash).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (r *userRepository) FindByUsername(username string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, username, password_hash, created_at::text
		FROM users
		WHERE username = $1
	`

	err := r.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with username %q not found: %w", username, sql.ErrNoRows)
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return user, nil
}

func (r *userRepository) FindByID(id int64) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, username, password_hash, created_at::text
		FROM users
		WHERE id = $1
	`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with id %d not found: %w", id, sql.ErrNoRows)
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return user, nil
}
