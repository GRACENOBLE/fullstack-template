package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"backend/internal/domain"
	"backend/internal/usecase"
)

type userRepository struct{ db *sql.DB }

// NewUserRepository returns a usecase.UserRepository backed by PostgreSQL.
func NewUserRepository(db *sql.DB) usecase.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Upsert(ctx context.Context, u *domain.User) (*domain.User, error) {
	const q = `
        INSERT INTO users (firebase_uid, name, email, photo_url)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (firebase_uid) DO UPDATE
            SET name       = EXCLUDED.name,
                email      = EXCLUDED.email,
                photo_url  = EXCLUDED.photo_url,
                updated_at = now()
        RETURNING id, firebase_uid, name, email, photo_url, created_at, updated_at`

	out := &domain.User{}
	err := r.db.QueryRowContext(ctx, q, u.FirebaseUID, u.Name, u.Email, u.PhotoURL).
		Scan(&out.ID, &out.FirebaseUID, &out.Name, &out.Email, &out.PhotoURL, &out.CreatedAt, &out.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("user repository: upsert: %w", err)
	}
	return out, nil
}

func (r *userRepository) DeleteByFirebaseUID(ctx context.Context, firebaseUID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE firebase_uid = $1`, firebaseUID)
	if err != nil {
		return fmt.Errorf("user repository: delete: %w", err)
	}
	return nil
}
