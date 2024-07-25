package dbrunner_test

import (
	"context"
	"testing"

	"github.com/database-playground/backend/internal/dbrunner"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunQuery(t *testing.T) {
	t.Parallel()

	t.Run("with valid query", func(t *testing.T) {
		t.Parallel()

		input := dbrunner.Input{
			Init: `
				CREATE TABLE test (
					id INTEGER PRIMARY KEY,
					name TEXT
				);

				INSERT INTO test (name) VALUES ('Alice');
				INSERT INTO test (name) VALUES ('Bob');
			`,
			Query: "SELECT * FROM test;",
		}
		output, err := dbrunner.RunQuery(context.Background(), input)
		require.NoError(t, err)
		assert.Equal(t, dbrunner.Output{
			Header: []string{"id", "name"},
			Data: [][]*string{
				{lo.ToPtr("1"), lo.ToPtr("Alice")},
				{lo.ToPtr("2"), lo.ToPtr("Bob")},
			},
		}, output)
	})

	t.Run("with no query", func(t *testing.T) {
		t.Parallel()

		input := dbrunner.Input{
			Init: `
				CREATE TABLE test (
					id INTEGER PRIMARY KEY,
					name TEXT
				);

				INSERT INTO test (name) VALUES ('Alice');
				INSERT INTO test (name) VALUES ('Bob');
			`,
			Query: "",
		}
		output, err := dbrunner.RunQuery(context.Background(), input)
		require.NoError(t, err)
		assert.Equal(t, dbrunner.Output{
			Header: []string{},
			Data:   [][]*string{},
		}, output)
	})

	t.Run("with UPDATE query", func(t *testing.T) {
		t.Parallel()

		input := dbrunner.Input{
			Init: `
				CREATE TABLE test (
					id INTEGER PRIMARY KEY,
					name TEXT
				);

				INSERT INTO test (name) VALUES ('Alice');
				INSERT INTO test (name) VALUES ('Bob');
			`,
			Query: "UPDATE test SET name = 'Charlie' WHERE id = 1;",
		}
		output, err := dbrunner.RunQuery(context.Background(), input)
		require.NoError(t, err)
		assert.Equal(t, dbrunner.Output{
			Header: []string{},
			Data:   [][]*string{},
		}, output)
	})

	t.Run("with UPDATE + RETURNING query", func(t *testing.T) {
		t.Parallel()

		input := dbrunner.Input{
			Init: `
				CREATE TABLE test (
					id INTEGER PRIMARY KEY,
					name TEXT
				);

				INSERT INTO test (name) VALUES ('Alice');
				INSERT INTO test (name) VALUES ('Bob');
			`,
			Query: "UPDATE test SET name = 'Charlie' WHERE id = 1 RETURNING *;",
		}
		output, err := dbrunner.RunQuery(context.Background(), input)
		require.NoError(t, err)
		assert.Equal(t, dbrunner.Output{
			Header: []string{"id", "name"},
			Data: [][]*string{
				{lo.ToPtr("1"), lo.ToPtr("Charlie")},
			},
		}, output)
	})

	t.Run("even with UPDATE query, the SELECT query should still stay as schema does", func(t *testing.T) {
		t.Parallel()

		initSQL := `
			CREATE TABLE test (
				id INTEGER PRIMARY KEY,
				name TEXT
			);

			INSERT INTO test (name) VALUES ('Alice');
			INSERT INTO test (name) VALUES ('Bob');
		`

		_, err := dbrunner.RunQuery(context.Background(), dbrunner.Input{
			Init:  initSQL,
			Query: "UPDATE test SET name = 'Charlie' WHERE id = 1;",
		})
		require.NoError(t, err)

		output, err := dbrunner.RunQuery(context.Background(), dbrunner.Input{
			Init:  initSQL,
			Query: "SELECT * FROM test;",
		})
		require.NoError(t, err)
		assert.Equal(t, dbrunner.Output{
			Header: []string{"id", "name"},
			Data: [][]*string{
				{lo.ToPtr("1"), lo.ToPtr("Alice")},
				{lo.ToPtr("2"), lo.ToPtr("Bob")},
			},
		}, output)
	})

	t.Run("with DoS query, it should be terminated", func(t *testing.T) {
		t.Parallel()

		input := dbrunner.Input{
			Init: `
				CREATE TABLE test (
					id INTEGER PRIMARY KEY,
					name TEXT
				);

				INSERT INTO test (name) VALUES ('Alice');
				INSERT INTO test (name) VALUES ('Bob');
			`,
			Query: `
				WITH RECURSIVE cte (n) AS (
					SELECT 1
					UNION ALL
					SELECT n + 1 FROM cte
				)
				SELECT * FROM cte;
			`,
		}

		_, err := dbrunner.RunQuery(context.Background(), input)
		require.Error(t, err)

		assert.ErrorIs(t, err, context.DeadlineExceeded)
	})

	t.Run("with malformed query, it should return an error", func(t *testing.T) {
		t.Parallel()

		input := dbrunner.Input{
			Init: `
				CREATE TABLE test (
					id INTEGER PRIMARY KEY,
					name TEXT
				);

				INSERT INTO test (name) VALUES ('Alice');
				INSERT INTO test (name) VALUES ('Bob');
			`,
			Query: "SELECT * FROM unknown_table;",
		}

		_, err := dbrunner.RunQuery(context.Background(), input)
		require.Error(t, err)

		assert.Contains(t, err.Error(), "no such table: unknown_table")
	})

	t.Run("with invalid query, it should return an error", func(t *testing.T) {
		t.Parallel()

		input := dbrunner.Input{
			Init: `
				CREATE TABLE test (
					id INTEGER PRIMARY KEY,
					name TEXT
				);

				INSERT INTO test (name) VALUES ('Alice');
				INSERT INTO test (name) VALUES ('Bob');
			`,
			Query: "SELECT * FROM test WHERE id = ':D)D)D))D)D)D)D)D;",
		}

		_, err := dbrunner.RunQuery(context.Background(), input)
		require.Error(t, err)
	})

	t.Run("with context cancelled, it should return an error", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		go cancel()

		input := dbrunner.Input{
			Init: `
				CREATE TABLE test (
					id INTEGER PRIMARY KEY,
					name TEXT
				);
				`,
			Query: "SELECT * FROM test;",
		}

		_, err := dbrunner.RunQuery(ctx, input)
		assert.ErrorIs(t, err, context.Canceled)
	})

	t.Run("with invalid schema, it should return an error", func(t *testing.T) {
		input := dbrunner.Input{
			Init: `
				CREATE TABLE test (
					id INTEGER PRIMARY KEY,
					name TEXT
				);

				INSERT INTO unknown_table (name) VALUES ('Alice');
				`,
			Query: "SELECT * FROM test;",
		}

		_, err := dbrunner.RunQuery(context.Background(), input)
		require.Error(t, err)

		assert.Contains(t, err.Error(), "exec init")
	})

	t.Run("with nil return, the cell should be <nil>", func(t *testing.T) {
		input := dbrunner.Input{
			Init: `
				CREATE TABLE test (
					id INTEGER PRIMARY KEY,
					name TEXT
				);

				INSERT INTO test VALUES (1, NULL);
				`,
			Query: "SELECT * FROM test;",
		}

		output, err := dbrunner.RunQuery(context.Background(), input)
		require.NoError(t, err)

		assert.Equal(t, dbrunner.Output{
			Header: []string{"id", "name"},
			Data: [][]*string{
				{lo.ToPtr("1"), nil},
			},
		}, output)
	})
}
