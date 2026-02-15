package models

import "time"

type RefreshSession struct {
	ID                  int64      `db:"id"`
	UserID              int64      `db:"user_id"`
	TokenHash           string     `db:"token_hash"`
	FamilyID            string     `db:"family_id"`
	IssuedAt            time.Time  `db:"issued_at"`
	ExpiresAt           time.Time  `db:"expires_at"`
	ConsumedAt          *time.Time `db:"consumed_at"`
	RevokedAt           *time.Time `db:"revoked_at"`
	ReplacedBySessionID *int64     `db:"replaced_by_session_id"`
}
