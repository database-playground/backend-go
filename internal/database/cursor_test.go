package database_test

import (
	"testing"

	"github.com/database-playground/backend/internal/database"
	"github.com/stretchr/testify/assert"
)

func TestCursor(t *testing.T) {
	t.Parallel()

	t.Run("limit 10, offset 0 by default", func(t *testing.T) {
		cursor := database.Cursor{}

		assert.Equal(t, 10, cursor.GetLimit())
		assert.Equal(t, 0, cursor.GetOffset())
	})

	t.Run("limit 99, offset 99 is acceptable", func(t *testing.T) {
		cursor := database.Cursor{Limit: 99, Offset: 99}

		assert.Equal(t, 99, cursor.GetLimit())
		assert.Equal(t, 99, cursor.GetOffset())
	})

	t.Run("limit over 100 should be rounded to 100", func(t *testing.T) {
		cursor := database.Cursor{Limit: 114514}

		assert.Equal(t, 100, cursor.GetLimit())
	})
}
