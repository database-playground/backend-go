package database_test

import (
	"context"
	"testing"
)

// fixme: manual comparing currently; need to mock the loggers

func TestMigration(t *testing.T) {
	db, cleanup := createOnetimeDatabase(t)
	defer cleanup()

	if err := db.Migrate(context.Background()); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
}

func TestDuplicateMigration(t *testing.T) {
	db, cleanup := createOnetimeDatabase(t)
	defer cleanup()

	if err := db.Migrate(context.Background()); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	if err := db.Migrate(context.Background()); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
}
