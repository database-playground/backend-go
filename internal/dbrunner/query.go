package dbrunner

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/samber/lo"
	_ "modernc.org/sqlite"
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
		values := make([]any, len(cols))
		for i := range values {
			values[i] = new(NullableStringScanner)
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
			row = append(row, struct {
				Column string
				Value  *string
			}{
				Column: col,
				Value:  values[i].(*NullableStringScanner).Value(),
			})
		}

		output.Result = append(output.Result, row)
	}
	if err := rows.Err(); err != nil {
		return Output{}, fmt.Errorf("rows error: %w", err)
	}

	return output, nil
}

type NullableStringScanner struct {
	value *string
}

func (n *NullableStringScanner) Scan(value any) error {
	if value == nil {
		n.value = nil
		return nil
	}

	n.value = lo.ToPtr(fmt.Sprintf("%v", value))
	return nil
}

func (n *NullableStringScanner) Value() *string {
	return n.value
}

var _ sql.Scanner = &NullableStringScanner{}
