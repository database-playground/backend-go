package database_test

import (
	"context"
	"testing"

	"github.com/database-playground/backend/internal/database"
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

	questions, err := db.ListQuestions(ctx, database.ListQuestionsParams{})
	if err != nil {
		t.Fatalf("failed to get schema: %v", err)
	}

	assert.Len(t, questions, 10) // default offset=0, limit=10
	assert.Equal(t, "Find a product in the shop", questions[0].Title)
	t.Logf("%#v", questions)
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

	// Get a question
	for _, listQuestion := range questions {
		getQuestion, err := db.GetQuestion(ctx, listQuestion.ID)
		if err != nil {
			t.Fatalf("failed to get schema: %v", err)
		}

		assert.Equal(t, listQuestion, getQuestion)
	}
}
