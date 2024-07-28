package database

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

var FxModule = fx.Module("database", fx.Provide(New), fx.Invoke(func(db *Database, lc fx.Lifecycle) {
	lc.Append(fx.Hook{
		OnStart: db.Migrate,
	})
}))

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

	return NewWithPool(pool, logger), nil
}

func NewWithPool(pool *pgxpool.Pool, logger *slog.Logger) *Database {
	return &Database{
		pool:   pool,
		logger: logger,
	}
}

func (db *Database) Close() {
	db.pool.Close()
}
