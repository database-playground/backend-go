package dbrunner_test

import (
	"testing"

	"github.com/database-playground/backend/internal/dbrunner"
)

func TestFormatSQL(t *testing.T) {
	t.Parallel()

	testmap := map[string]string{
		"SELECT * FROM table":                                 "SELECT * FROM TABLE",
		"SELECT * FROM table WHERE id = 1":                    "SELECT * FROM TABLE WHERE id = 1",
		"SELECT * FROM table WHERE id = 1;":                   "SELECT * FROM TABLE WHERE id = 1",
		"SELECT * FROM table WHERE id = 1; -- comment":        "SELECT * FROM TABLE WHERE id = 1",
		"SELECT * FROM table WHERE id = 1; -- comment\n":      "SELECT * FROM TABLE WHERE id = 1",
		"SELECT *, aaa FROM table WHERE id = 1; -- comment\n": "SELECT *, aaa FROM TABLE WHERE id = 1",
		"SELECT * FROM table;\nSELECT * FROM abc;":            "SELECT * FROM TABLE; SELECT * FROM abc",
		"SELECT *     FROM   table":                           "SELECT * FROM TABLE",
		"seLect * fRom table":                                 "SELECT * FROM TABLE",
	}

	for raw, expected := range testmap {
		t.Run(raw, func(t *testing.T) {
			t.Parallel()

			normalized, err := dbrunner.FormatSQL(raw)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if normalized != expected {
				t.Errorf("expected: %s, got: %s", expected, normalized)
			}
		})
	}
}
