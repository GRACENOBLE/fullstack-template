package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"backend/internal/handler"
	"backend/internal/repository/postgres"
	"backend/internal/usecase"
)

// NewServer wires all layers and returns a configured *http.Server.
func NewServer() (*http.Server, error) {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	cfg := postgres.DBConfig{
		Host:     os.Getenv("BLUEPRINT_DB_HOST"),
		Port:     os.Getenv("BLUEPRINT_DB_PORT"),
		Database: os.Getenv("BLUEPRINT_DB_DATABASE"),
		Username: os.Getenv("BLUEPRINT_DB_USERNAME"),
		Password: os.Getenv("BLUEPRINT_DB_PASSWORD"),
		Schema:   os.Getenv("BLUEPRINT_DB_SCHEMA"),
		SSLMode:  os.Getenv("BLUEPRINT_DB_SSLMODE"),
	}

	db, err := postgres.NewPostgresDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("server: database: %w", err)
	}

	healthRepo := postgres.NewHealthRepository(db)
	healthUC := usecase.NewHealthUseCase(healthRepo)
	h := handler.NewHandler(healthUC)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      h.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return srv, nil
}
