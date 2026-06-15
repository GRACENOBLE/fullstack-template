package handlers

import (
	"backend/internal/domain"
	"backend/internal/usecase"
)

// HealthStats is a type alias so swaggo can resolve domain.HealthStats
// from within this package for annotation purposes.
type HealthStats = domain.HealthStats

// FirebaseToken is a type alias so swaggo can resolve usecase.FirebaseToken
// from within this package for annotation purposes.
type FirebaseToken = usecase.FirebaseToken
