package storage

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func MustLoadDatabase(databaseURL string) *pgxpool.Pool {
	ctx := context.Background()

	dbpool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := dbpool.Ping(pingCtx); err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	db := stdlib.OpenDBFromPool(dbpool)

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("Failed to set goose dialect: %v", err)
	}

	log.Println("Running migrations...")
	if err := goose.Up(db, "migrations"); err != nil {
		log.Fatalf("Goose migrations failed: %v", err)
	}

	log.Println("Successfully connected to database and applied migrations!")
	return dbpool
}
