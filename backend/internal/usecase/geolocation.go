package usecase

import (
	"context"

	"backend/internal/domain"
)

// GeoLocator resolves an IP address to geographic metadata.
type GeoLocator interface {
	Lookup(ctx context.Context, ip string) (*domain.GeoLocation, error)
}
