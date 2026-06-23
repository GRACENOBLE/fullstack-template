package usecase

import (
	"context"

	"backend/internal/domain"
)

// UserRepository persists user profile data.
type UserRepository interface {
	Upsert(ctx context.Context, u *domain.User) (*domain.User, error)
	DeleteByFirebaseUID(ctx context.Context, firebaseUID string) error
}
