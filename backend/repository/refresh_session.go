package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"goTodo/backend/models"
)

type RefreshSessionRepository interface {
	CreateSession(userID int64, familyID string, tokenHash string, expiresAt time.Time) (int64, error)
	FindByTokenHash(tokenHash string) (*models.RefreshSession, error)
	RotateSession(oldSessionID int64, userID int64, familyID string, newTokenHash string, newExpiresAt time.Time) error
	RevokeFamily(familyID string, reason string) error
	RevokeByTokenHash(tokenHash string, reason string) error
}

type refreshSessionRepository struct {
	db *sql.DB
}

func NewRefreshSessionRepository(db *sql.DB) RefreshSessionRepository {
	return &refreshSessionRepository{db: db}
}

func (r *refreshSessionRepository) CreateSession(userID int64, familyID string, tokenHash string, expiresAt time.Time) (int64, error) {
	var sessionID int64
	query := `
		INSERT INTO auth_refresh_sessions (user_id, token_hash, family_id, issued_at, expires_at)
		VALUES ($1, $2, $3::uuid, NOW(), $4)
		RETURNING id
	`

	err := r.db.QueryRow(query, userID, tokenHash, familyID, expiresAt).Scan(&sessionID)
	if err != nil {
		return 0, fmt.Errorf("failed to create refresh session: %w", err)
	}

	return sessionID, nil
}

func (r *refreshSessionRepository) FindByTokenHash(tokenHash string) (*models.RefreshSession, error) {
	session := &models.RefreshSession{}
	query := `
		SELECT id, user_id, token_hash, family_id::text, issued_at, expires_at, consumed_at, revoked_at, replaced_by_session_id
		FROM auth_refresh_sessions
		WHERE token_hash = $1
	`

	err := r.db.QueryRow(query, tokenHash).Scan(
		&session.ID,
		&session.UserID,
		&session.TokenHash,
		&session.FamilyID,
		&session.IssuedAt,
		&session.ExpiresAt,
		&session.ConsumedAt,
		&session.RevokedAt,
		&session.ReplacedBySessionID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to find refresh session by token hash: %w", err)
	}

	return session, nil
}

func (r *refreshSessionRepository) RotateSession(oldSessionID int64, userID int64, familyID string, newTokenHash string, newExpiresAt time.Time) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin refresh rotation transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var newSessionID int64
	insertQuery := `
		INSERT INTO auth_refresh_sessions (user_id, token_hash, family_id, issued_at, expires_at)
		VALUES ($1, $2, $3::uuid, NOW(), $4)
		RETURNING id
	`
	if err = tx.QueryRow(insertQuery, userID, newTokenHash, familyID, newExpiresAt).Scan(&newSessionID); err != nil {
		return fmt.Errorf("failed to insert rotated refresh session: %w", err)
	}

	updateQuery := `
		UPDATE auth_refresh_sessions
		SET consumed_at = NOW(), replaced_by_session_id = $2, updated_at = NOW()
		WHERE id = $1
	`
	if _, err = tx.Exec(updateQuery, oldSessionID, newSessionID); err != nil {
		return fmt.Errorf("failed to mark old refresh session as consumed: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit refresh rotation transaction: %w", err)
	}

	return nil
}

func (r *refreshSessionRepository) RevokeFamily(familyID string, reason string) error {
	query := `
		UPDATE auth_refresh_sessions
		SET revoked_at = NOW(), revoke_reason = $2, updated_at = NOW()
		WHERE family_id = $1::uuid AND revoked_at IS NULL
	`

	if _, err := r.db.Exec(query, familyID, reason); err != nil {
		return fmt.Errorf("failed to revoke refresh family: %w", err)
	}

	return nil
}

func (r *refreshSessionRepository) RevokeByTokenHash(tokenHash string, reason string) error {
	query := `
		UPDATE auth_refresh_sessions
		SET revoked_at = NOW(), revoke_reason = $2, updated_at = NOW()
		WHERE token_hash = $1 AND revoked_at IS NULL
	`

	if _, err := r.db.Exec(query, tokenHash, reason); err != nil {
		return fmt.Errorf("failed to revoke refresh session by token hash: %w", err)
	}

	return nil
}
