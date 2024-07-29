//go:generate go run github.com/jmattheis/goverter/cmd/goverter gen .

package gatewayservice

import (
	"strconv"
	"time"

	"github.com/database-playground/backend/internal/models"
)

// goverter:converter
// goverter:matchIgnoreCase
// goverter:extend Int64ToString
// goverter:extend PInt64ToPString
// goverter:extend TimeToPtrTime
type Converter interface {
	SchemaFromModel(in *models.Schema) *Schema
	SchemaInitialSQLFromModel(in *models.SchemaInitialSQL) *SchemaInitialSQL
	// goverter:enum:unknown Empty
	// goverter:enum:map DifficultyEasy Easy
	// goverter:enum:map DifficultyMedium Medium
	// goverter:enum:map DifficultyHard Hard
	DifficultyFromModel(in models.Difficulty) QuestionDifficulty
	QuestionFromModel(in *models.Question) *Question
	QuestionAnswerFromModel(in *models.QuestionAnswer) *QuestionAnswer
	QuestionSolutionFromModel(in *models.QuestionSolution) *QuestionSolution
}

func Int64ToString(in int64) string {
	return strconv.FormatInt(in, 10)
}

func PInt64ToPString(in *int64) *string {
	if in == nil {
		return nil
	}
	out := Int64ToString(*in)
	return &out
}

func TimeToPtrTime(in time.Time) *time.Time {
	return &in
}
