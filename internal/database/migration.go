package database

import (
	"context"
	"embed"
	"fmt"
	"log/slog"
	"path/filepath"
	"slices"
)

//go:embed migrations
var migrationsFs embed.FS

//go:embed seeds
var seedsFs embed.FS

var availableMigrations []migrationFile

func init() {
	migrations, err := migrationsFs.ReadDir("migrations")
	if err != nil {
		panic(fmt.Errorf("read migrations directory: %w", err))
	}

	availableMigrations = make([]migrationFile, 0, len(migrations))

	for _, migration := range migrations {
		migrationf := migrationFile{
			Name: migration.Name(),
		}

		if !migration.Type().IsRegular() {
			continue
		}

		content, err := migrationsFs.ReadFile(filepath.Join("migrations", migration.Name()))
		if err != nil {
			panic(fmt.Errorf("read migration file: %w", err))
		}

		migrationf.Content = string(content)

		availableMigrations = append(availableMigrations, migrationf)
	}
}

type migrationFile struct {
	Name    string
	Content string
}

// Migrate runs the database migrations according to the version
func (db *Database) Migrate(ctx context.Context) error {
	conn, err := db.pool.Acquire(ctx)
	if err != nil {
		db.logger.Error("failed to acquire connection from pool", slog.Any("error", err))
		return fmt.Errorf("acquire connection from pool: %w", err)
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		db.logger.Error("failed to begin transaction", slog.Any("error", err))
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	// Create migration table to track how many migrations have been run
	_, err = tx.Exec(ctx, `
		--sql
		CREATE TABLE IF NOT EXISTS dp_migrations (
			migration_version VARCHAR(255) PRIMARY KEY
		);
	`)
	if err != nil {
		db.logger.Error("failed to create migrations table", slog.Any("error", err))
		return fmt.Errorf("create migrations table: %w", err)
	}

	// Get what migrations have been run
	var ranMigrations []string
	rows, err := tx.Query(ctx, "SELECT migration_version FROM dp_migrations")
	if err != nil {
		db.logger.Error("failed to get ran migrations", slog.Any("error", err))
	}
	for rows.Next() {
		var migration string
		err = rows.Scan(&migration)
		if err != nil {
			break
		}
		ranMigrations = append(ranMigrations, migration)
	}
	if rows.Err() != nil {
		db.logger.Error("failed to scan ran migrations", slog.Any("error", rows.Err()))
		return fmt.Errorf("scan ran migrations: %w", rows.Err())
	}

	for _, migration := range availableMigrations {
		if slices.Contains(ranMigrations, migration.Name) {
			db.logger.Info("migration already ran", slog.Any("version", migration.Name))
			continue
		}

		db.logger.Info("migrating database", slog.Any("version", migration.Name))

		_, err = tx.Exec(ctx, migration.Content)
		if err != nil {
			db.logger.Error("failed to run migration", slog.Any("error", err))
			return fmt.Errorf("run migration: %w", err)
		}

		// Insert the migration into the table
		_, err = tx.Exec(ctx, "INSERT INTO dp_migrations (migration_version) VALUES ($1)", migration.Name)
		if err != nil {
			db.logger.Error("failed to insert migration into table", slog.Any("error", err))
			return fmt.Errorf("insert migration into table: %w", err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		db.logger.Error("failed to commit transaction", slog.Any("error", err))
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

// SeedTestOnly seeds the database with test data.
//
// It must not be used in production since it does not check for duplicates.
// It is intended to be used only in tests.
func (db *Database) SeedTestOnly(ctx context.Context) error {
	seedsFiles, err := seedsFs.ReadDir("seeds")
	if err != nil {
		db.logger.Error("failed to read seeds directory", slog.Any("error", err))
		return fmt.Errorf("read seeds directory: %w", err)
	}

	for _, seed := range seedsFiles {
		if !seed.Type().IsRegular() {
			continue
		}

		content, err := seedsFs.ReadFile(filepath.Join("seeds", seed.Name()))
		if err != nil {
			db.logger.Error("failed to read seed file", slog.Any("error", err))
			return fmt.Errorf("read seed file: %w", err)
		}

		_, err = db.pool.Exec(ctx, string(content))
		if err != nil {
			db.logger.Error("failed to seed database", slog.Any("error", err))
			return fmt.Errorf("seed database: %w", err)
		}

		db.logger.Info("seeded database", slog.Any("file", seed.Name()))
	}

	return nil
}
