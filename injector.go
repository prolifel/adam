package main

import (
	"database/sql"
	"fmt"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
	"github.com/samber/do/v2"
)

func initDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./container_profiles.db")
	if err != nil {
		return nil, err
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
	if err != nil {
		panic(fmt.Errorf("Failed to get .env file: %+v", err))
	}

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
