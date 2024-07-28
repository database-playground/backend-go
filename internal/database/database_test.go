package database_test

import (
	"context"
	"log/slog"
	"math/rand/v2"
	"os"
	"strconv"
	"testing"

	"github.com/database-playground/backend/internal/database"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func createOnetimeDatabase(t *testing.T) (*database.Database, func()) {
	uri := os.Getenv("TEST_POSTGRES_URI")
	if uri == "" {
		t.Skip("skipping test; no database connection")
	}

	rootConnection, err := pgx.Connect(context.Background(), uri)
	if err != nil {
		t.Skip("skipping test; unconnectable database")
	}

	newDatabase := "dp_backend_test_" + strconv.FormatInt(rand.Int64(), 36)
	_, err = rootConnection.Exec(context.Background(), "CREATE DATABASE "+newDatabase)
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}
	t.Logf("database %s created", newDatabase)

	newConnConfig, err := pgxpool.ParseConfig(uri)
	if err != nil {
		t.Fatalf("failed to parse connection string: %v", err)
	}
	newConnConfig.ConnConfig.Database = newDatabase

	newPool, err := pgxpool.NewWithConfig(context.Background(), newConnConfig)
	if err != nil {
		t.Fatalf("failed to create pool: %v", err)
	}

	return database.NewWithPool(newPool, slog.Default()), func() {
		// close pool
		newPool.Close()

		// use the rootconnection to drop the database
		_, err := rootConnection.Exec(context.Background(), "DROP DATABASE "+newDatabase)
		if err != nil {
			t.Errorf("failed to drop database: %v", err)
		} else {
			t.Logf("database %s dropped", newDatabase)
		}

		// close root connection
		if err := rootConnection.Close(context.Background()); err != nil {
			t.Fatalf("failed to close connection: %v", err)
		}
	}
}
