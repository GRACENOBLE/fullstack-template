package postgres

import (
	"context"
	"database/sql"
	"log"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testDB *sql.DB

func mustStartPostgresContainer() (func(context.Context, ...testcontainers.TerminateOption) error, error) {
	var (
		dbName = "database"
		dbPwd  = "password"
		dbUser = "user"
	)

	container, err := tcpostgres.Run(
		context.Background(),
		"postgres:latest",
		tcpostgres.WithDatabase(dbName),
		tcpostgres.WithUsername(dbUser),
		tcpostgres.WithPassword(dbPwd),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)
	if err != nil {
		return nil, err
	}

	dbHost, err := container.Host(context.Background())
	if err != nil {
		return container.Terminate, err
	}

	dbPort, err := container.MappedPort(context.Background(), "5432/tcp")
	if err != nil {
		return container.Terminate, err
	}

	cfg := DBConfig{
		Host:     dbHost,
		Port:     dbPort.Port(),
		Database: dbName,
		Username: dbUser,
		Password: dbPwd,
		Schema:   "public",
		SSLMode:  "disable",
	}

	db, err := NewPostgresDB(cfg)
	if err != nil {
		return container.Terminate, err
	}
	testDB = db

	return container.Terminate, nil
}

func TestMain(m *testing.M) {
	teardown, err := mustStartPostgresContainer()
	if err != nil {
		log.Fatalf("could not start postgres container: %v", err)
	}

	m.Run()

	if teardown != nil && teardown(context.Background()) != nil {
		log.Fatalf("could not teardown postgres container: %v", err)
	}
}

func TestNew(t *testing.T) {
	if testDB == nil {
		t.Fatal("NewPostgresDB() returned nil db")
	}
}

func TestHealth(t *testing.T) {
	repo := NewHealthRepository(testDB)

	stats, err := repo.Health(context.Background())
	if err != nil {
		t.Fatalf("Health() returned unexpected error: %v", err)
	}
	if stats.Status != "up" {
		t.Fatalf("expected status to be up, got %s", stats.Status)
	}
	if stats.Error != "" {
		t.Fatalf("expected no error, got %s", stats.Error)
	}
}
