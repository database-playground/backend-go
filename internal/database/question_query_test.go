package database_test

import (
	"context"
	"testing"

	"github.com/database-playground/backend/internal/database"
	"github.com/database-playground/backend/internal/models"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListQuestions(t *testing.T) {
	t.Parallel()

	db, cleanup := createOnetimeDatabase(t)
	defer cleanup()

	ctx := context.Background()
	require.NoError(t, db.Migrate(ctx))
	require.NoError(t, db.SeedTestOnly(ctx))

	rootQuestions, err := db.ListQuestions(ctx, database.ListQuestionsParams{})
	if err != nil {
		t.Fatalf("failed to get schema: %v", err)
	}

	assert.Len(t, rootQuestions, 10) // default offset=0, limit=10
	assert.Equal(t, "Find a product in the shop", rootQuestions[0].Title)
	t.Logf("%#v", rootQuestions)

	t.Run("offset=1; limit=5", func(t *testing.T) {
		questions, err := db.ListQuestions(ctx, database.ListQuestionsParams{Cursor: database.Cursor{Offset: 1, Limit: 5}})
		if err != nil {
			t.Fatalf("failed to get schema: %v", err)
		}

		t.Logf("%#v", lo.Map(questions, func(q *models.Question, _ int) string {
			return q.Title
		}))

		assert.Len(t, questions, 5)

		for i, q := range questions {
			assert.Equal(t, rootQuestions[i+1].Title, q.Title)
		}
	})

	t.Run("offset=5; limit=5", func(t *testing.T) {
		questions, err := db.ListQuestions(ctx, database.ListQuestionsParams{Cursor: database.Cursor{Offset: 5, Limit: 5}})
		if err != nil {
			t.Fatalf("failed to get schema: %v", err)
		}

		assert.Len(t, questions, 5)

		for i, q := range questions {
			assert.Equal(t, rootQuestions[i+5].Title, q.Title)
		}
	})

	t.Run("offset=0; limit=5", func(t *testing.T) {
		questions, err := db.ListQuestions(ctx, database.ListQuestionsParams{Cursor: database.Cursor{Offset: 0, Limit: 5}})
		if err != nil {
			t.Fatalf("failed to get schema: %v", err)
		}

		assert.Len(t, questions, 5)

		for i, q := range questions {
			assert.Equal(t, rootQuestions[i].Title, q.Title)
		}
	})

	t.Run("limit=0 should be represented as limit=10", func(t *testing.T) {
		questions, err := db.ListQuestions(ctx, database.ListQuestionsParams{Cursor: database.Cursor{Limit: 0}})
		if err != nil {
			t.Fatalf("failed to get schema: %v", err)
		}

		assert.Len(t, questions, 10)

		for i, q := range questions {
			assert.Equal(t, rootQuestions[i].Title, q.Title)
		}
	})
}

func TestGetQuestion(t *testing.T) {
	t.Parallel()

	db, cleanup := createOnetimeDatabase(t)
	defer cleanup()

	ctx := context.Background()
	require.NoError(t, db.Migrate(ctx))
	require.NoError(t, db.SeedTestOnly(ctx))

	// List questions
	questions, err := db.ListQuestions(ctx, database.ListQuestionsParams{})
	if err != nil {
		t.Fatalf("failed to get schema: %v", err)
	}

	// Check if GetQuestion returns the same question
	for _, listQuestion := range questions {
		getQuestion, err := db.GetQuestion(ctx, listQuestion.ID)
		if err != nil {
			t.Fatalf("failed to get schema: %v", err)
		}

		assert.Equal(t, listQuestion, getQuestion)
	}
}
