package database_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/database-playground/backend/internal/database"
)

// fixme: this test should not require the active connection :(

func TestMigration(t *testing.T) {
	testPostgresUri := os.Getenv("MIGRATION_TEST_POSTGRES_URI")
	if testPostgresUri == "" {
		t.Skip("skipping test; no database connection")
	}

	database, err := database.NewWithURI(testPostgresUri, slog.Default())
	if err != nil {
		t.Skip("skipping test; unconnectable database")
	}

	if err := database.Migrate(context.Background()); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
}
