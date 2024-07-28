package database_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetSchema(t *testing.T) {
	t.Parallel()

	db, cleanup := createOnetimeDatabase(t)
	defer cleanup()

	ctx := context.Background()
	require.NoError(t, db.Migrate(ctx))
	require.NoError(t, db.SeedTestOnly(ctx))

	schema, err := db.GetSchema(ctx, "shop")
	if err != nil {
		t.Fatalf("failed to get schema: %v", err)
	}

	assert.Equal(t, "shop", schema.SchemaID)
	assert.Equal(t, "The schema that is for a shop", schema.Description)

	t.Logf("%#v", schema)
}

func TestGetSchemaInitialSQL(t *testing.T) {
	t.Parallel()

	db, cleanup := createOnetimeDatabase(t)
	defer cleanup()

	ctx := context.Background()
	require.NoError(t, db.Migrate(ctx))
	require.NoError(t, db.SeedTestOnly(ctx))

	sql, err := db.GetSchemaInitialSQL(ctx, "shop")
	if err != nil {
		t.Fatalf("failed to get schema initial sql: %v", err)
	}

	assert.Contains(t, sql, "CREATE TABLE products")
}
