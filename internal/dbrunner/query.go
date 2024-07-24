package dbrunner

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/samber/lo"
	"modernc.org/sqlite"
)

const timeoutSecond = 1

func RunQuery(ctx context.Context, input Input) (Output, error) {
	ctx, cancel := context.WithTimeout(ctx, timeoutSecond*time.Second)
	defer cancel()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return Output{}, fmt.Errorf("open database: %w", err)
	}

	_, err = db.ExecContext(ctx, input.Init)
	if err != nil {
		return Output{}, fmt.Errorf("exec init: %w", err)
	}

	rows, err := db.QueryContext(ctx, input.Query)
	if err != nil {
		return Output{}, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return Output{}, fmt.Errorf("get columns: %w", err)
	}

	output := Output{}
	for rows.Next() {
		var values []any
		for range cols {
			var value any
			values = append(values, &value)
		}

		err := rows.Scan(values...)
		if err != nil {
			break
		}

		var row []struct {
			Column string
			Value  *string
		}
		for i, col := range cols {
			var value *string
			if values[i] != nil {
				value = lo.ToPtr(fmt.Sprint(*values[i].(*any)))
			} else {
				value = nil
			}

			row = append(row, struct {
				Column string
				Value  *string
			}{
				Column: col,
				Value:  value,
			})
		}

		output.Result = append(output.Result, row)
	}
	if err := rows.Err(); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return Output{}, fmt.Errorf("query timeout: %w", err)
		}

		var sqliteError *sqlite.Error
		if errors.As(err, &sqliteError) {
			return Output{}, fmt.Errorf("sqlite error: %w", err)
		}

		return Output{}, fmt.Errorf("rows error: %w", err)
	}

	return output, nil
}
