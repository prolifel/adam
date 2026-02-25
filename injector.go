package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
	"github.com/samber/do/v2"
)

func initDB() (*sql.DB, error) {
	// Get database path from environment or use default
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./container_profiles.db"
	}

	// Ensure the directory exists
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Run migrations
	if err := goose.SetDialect("sqlite3"); err != nil {
		return nil, err
	}

	if err := goose.Up(db, "./migrations"); err != nil {
		return nil, err
	}

	return db, nil
}

func startProgram() do.Injector {
	err := godotenv.Load()

	var cfg Config
	err = env.Parse(&cfg)
	if err != nil {
		panic(fmt.Errorf("Failed to parse .env file: %+v", err))
	}

	db, err := initDB()
	if err != nil {
		panic(fmt.Errorf("Failed to initialize database: %+v", err))
	}

	injector := do.New()

	// Provide config
	do.ProvideValue(injector, cfg)

	// Provide database
	do.ProvideValue(injector, db)

	// Provide Repository
	do.Provide(injector, func(i do.Injector) (*Repo, error) {
		repo := &Repo{
			DB: do.MustInvoke[*sql.DB](i),
		}
		return repo, nil
	})

	// Provide Service
	do.Provide(injector, func(i do.Injector) (*Service, error) {
		service := &Service{
			Repo: do.MustInvoke[*Repo](i),
			Cfg:  do.MustInvoke[Config](i),
		}
		return service, nil
	})

	return injector
}
