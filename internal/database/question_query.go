package database

import (
	"context"

	"github.com/database-playground/backend/internal/models"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/samber/lo"
)

type ListQuestionsParams struct {
	Offset int
	Limit  int
}

func (db *Database) ListQuestions(ctx context.Context, param ListQuestionsParams) ([]*models.Question, error) {
	var questions []*models.Question

	limit := lo.CoalesceOrEmpty(param.Limit, 10)

	err := pgxscan.Select(ctx, db.pool, &questions, `
		--sql
		SELECT question_id, schema_id, type, difficulty, title, description, created_at, updated_at
		FROM dp_questions
		ORDER BY question_id
		LIMIT $1 OFFSET $2;
	`, limit, param.Offset)
	if err != nil {
		return nil, err
	}

	return questions, nil
}

func (db *Database) GetQuestion(ctx context.Context, questionID int64) (*models.Question, error) {
	var question models.Question

	err := pgxscan.Get(ctx, db.pool, &question, `
		--sql
		SELECT question_id, schema_id, type, difficulty, title, description, created_at, updated_at
		FROM dp_questions
		WHERE question_id = $1;
	`, questionID)
	if err != nil {
		return nil, err
	}

	return &question, nil
}
