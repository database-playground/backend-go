package database

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

func New(logger *slog.Logger) (*Database, error) {
	postgresURI := os.Getenv("POSTGRES_URI")
	if postgresURI == "" {
		return nil, errors.New("POSTGRES_URI is required")
	}

	return NewWithURI(postgresURI, logger)
}

func NewWithURI(postgresURI string, logger *slog.Logger) (*Database, error) {
	pool, err := pgxpool.New(context.Background(), postgresURI)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	return &Database{
		pool:   pool,
		logger: logger,
	}, nil
}

func (db *Database) Close() {
	db.pool.Close()
}
