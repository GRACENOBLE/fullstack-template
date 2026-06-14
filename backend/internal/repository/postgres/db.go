package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// DBConfig holds all parameters needed to open a PostgreSQL connection.
type DBConfig struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
	Schema   string
	SSLMode  string
}

// NewPostgresDB opens a *sql.DB connection using the provided config.
// The caller is responsible for calling db.Close() on shutdown.
func NewPostgresDB(cfg DBConfig) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s&search_path=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.SSLMode, cfg.Schema,
	)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("postgres: open: %w", err)
	}
	return db, nil
}
