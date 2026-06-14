package usecase

import (
	"context"

	"backend/internal/domain"
)

// HealthReader is the repository interface the use case depends on.
type HealthReader interface {
	Health(ctx context.Context) (domain.HealthStats, error)
}

// HealthUseCase is the application-layer interface consumed by handlers.
type HealthUseCase interface {
	GetHealth(ctx context.Context) (domain.HealthStats, error)
}

type healthUseCase struct {
	repo HealthReader
}

// NewHealthUseCase constructs the use case with its repository dependency.
func NewHealthUseCase(repo HealthReader) HealthUseCase {
	return &healthUseCase{repo: repo}
}

func (uc *healthUseCase) GetHealth(ctx context.Context) (domain.HealthStats, error) {
	stats, err := uc.repo.Health(ctx)
	if err != nil {
		return stats, err
	}

	stats.Message = "It's healthy"
	if stats.OpenConnections > 40 {
		stats.Message = "The database is experiencing heavy load."
	}
	if stats.WaitCount > 1000 {
		stats.Message = "The database has a high number of wait events, indicating potential bottlenecks."
	}
	if stats.MaxIdleClosed > int64(stats.OpenConnections)/2 {
		stats.Message = "Many idle connections are being closed, consider revising the connection pool settings."
	}
	if stats.MaxLifetimeClosed > int64(stats.OpenConnections)/2 {
		stats.Message = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats, nil
}
