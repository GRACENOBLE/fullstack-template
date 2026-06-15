package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"backend/internal/domain"
)

// FCMTokenRepository implements usecase.FCMTokenRepository against PostgreSQL.
type FCMTokenRepository struct{ db *sql.DB }

// NewFCMTokenRepository creates a new FCMTokenRepository.
func NewFCMTokenRepository(db *sql.DB) *FCMTokenRepository {
	return &FCMTokenRepository{db: db}
}

// SaveToken persists a device token for a user. If the token already exists it is
// updated (upserted) so that a single physical device is always associated with
// exactly one user.
func (r *FCMTokenRepository) SaveToken(ctx context.Context, userID, token, platform string) error {
	const q = `
		INSERT INTO fcm_tokens (user_id, token, platform)
		VALUES ($1, $2, $3)
		ON CONFLICT (token) DO UPDATE
			SET user_id  = EXCLUDED.user_id,
			    platform = EXCLUDED.platform`
	if _, err := r.db.ExecContext(ctx, q, userID, token, platform); err != nil {
		return fmt.Errorf("fcm_token_repository: save: %w", err)
	}
	return nil
}

// GetTokensByUserID returns all FCM tokens registered for a user.
func (r *FCMTokenRepository) GetTokensByUserID(ctx context.Context, userID string) ([]domain.FCMToken, error) {
	const q = `SELECT id, user_id, token, platform, created_at FROM fcm_tokens WHERE user_id = $1`
	rows, err := r.db.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("fcm_token_repository: get: %w", err)
	}
	defer rows.Close()

	var tokens []domain.FCMToken
	for rows.Next() {
		var t domain.FCMToken
		if err := rows.Scan(&t.ID, &t.UserID, &t.Token, &t.Platform, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("fcm_token_repository: scan: %w", err)
		}
		tokens = append(tokens, t)
	}
	return tokens, rows.Err()
}

// DeleteToken removes a single FCM token, typically called on logout.
func (r *FCMTokenRepository) DeleteToken(ctx context.Context, token string) error {
	const q = `DELETE FROM fcm_tokens WHERE token = $1`
	if _, err := r.db.ExecContext(ctx, q, token); err != nil {
		return fmt.Errorf("fcm_token_repository: delete: %w", err)
	}
	return nil
}
