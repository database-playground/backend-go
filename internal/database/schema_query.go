package database

import (
	"context"

	"github.com/database-playground/backend/internal/models"
	"github.com/georgysavva/scany/v2/pgxscan"
)

func (db *Database) GetSchema(ctx context.Context, schemaID string) (*models.Schema, error) {
	var schema models.Schema

	err := pgxscan.Get(ctx, db.pool, &schema, `
		--sql
		SELECT schema_id, picture, description, created_at, updated_at
		FROM dp_schemas
		WHERE schema_id = $1;
	`, schemaID)
	if err != nil {
		return nil, err
	}

	return &schema, nil
}

func (db *Database) GetSchemaInitialSQL(ctx context.Context, schemaID string) (*models.SchemaInitialSQL, error) {
	var model models.SchemaInitialSQL

	err := pgxscan.Get(ctx, db.pool, &model, `
		--sql
		SELECT schema_id, initial_sql
		FROM dp_schemas
		WHERE schema_id = $1;
	`, schemaID)
	if err != nil {
		return nil, err
	}

	return &model, nil
}
