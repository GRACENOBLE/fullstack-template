package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"backend/internal/domain"
)

// HealthRepository checks the health of the postgres connection.
type HealthRepository struct {
	db *sql.DB
}

// NewHealthRepository constructs a HealthRepository.
func NewHealthRepository(db *sql.DB) *HealthRepository {
	return &HealthRepository{db: db}
}

// Health pings the database and returns raw connection pool statistics.
// Returns stats with status "down" and a non-nil error when the ping fails.
func (r *HealthRepository) Health(ctx context.Context) (domain.HealthStats, error) {
	pingCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if err := r.db.PingContext(pingCtx); err != nil {
		return domain.HealthStats{
			Status: "down",
			Error:  fmt.Sprintf("db down: %v", err),
		}, fmt.Errorf("postgres: health ping: %w", err)
	}

	dbStats := r.db.Stats()
	return domain.HealthStats{
		Status:            "up",
		OpenConnections:   dbStats.OpenConnections,
		InUse:             dbStats.InUse,
		Idle:              dbStats.Idle,
		WaitCount:         dbStats.WaitCount,
		WaitDuration:      dbStats.WaitDuration.String(),
		MaxIdleClosed:     dbStats.MaxIdleClosed,
		MaxLifetimeClosed: dbStats.MaxLifetimeClosed,
	}, nil
}
