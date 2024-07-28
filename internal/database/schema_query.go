package database

import (
	"context"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
)

type GetSchemaModel struct {
	SchemaID    string
	Picture     *string
	Description string

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (db *Database) GetSchema(ctx context.Context, schemaID string) (GetSchemaModel, error) {
	var schema GetSchemaModel

	err := pgxscan.Get(ctx, db.pool, &schema, `
		--sql
		SELECT schema_id, picture, description, created_at, updated_at
		FROM dp_schemas
		WHERE schema_id = $1;
	`, schemaID)
	if err != nil {
		return schema, err
	}

	return schema, nil
}

func (db *Database) GetSchemaInitialSQL(ctx context.Context, schemaID string) (string, error) {
	var schema string

	err := pgxscan.Get(ctx, db.pool, &schema, `
		--sql
		SELECT initial_sql
		FROM dp_schemas
		WHERE schema_id = $1;
	`, schemaID)
	if err != nil {
		return schema, err
	}

	return schema, nil
}
