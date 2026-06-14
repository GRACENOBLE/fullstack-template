package main

import (
	"context"
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/pressly/goose/v3"

	"backend/internal/infrastructure/database/postgres"
)

const migrationsDir = "internal/infrastructure/database/migrations"

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: migrate <command> [args]")
		fmt.Fprintln(os.Stderr, "commands: up, up-by-one, down, down-to <version>, reset, status, version, create <name>")
		os.Exit(1)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("migrate: set dialect: %v", err)
	}

	command, args := os.Args[1], os.Args[2:]

	// create only touches the filesystem — connect to DB is not needed
	if command == "create" {
		if len(args) == 0 {
			log.Fatal("migrate create: name is required (e.g. make migrate-create name=add_users)")
		}
		if err := goose.Create(nil, migrationsDir, args[0], "sql"); err != nil {
			log.Fatalf("migrate create: %v", err)
		}
		return
	}

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
		log.Fatalf("migrate: database: %v", err)
	}
	defer db.Close()

	if err := goose.RunContext(context.Background(), command, db, migrationsDir, args...); err != nil {
		log.Fatalf("migrate %s: %v", command, err)
	}
}
