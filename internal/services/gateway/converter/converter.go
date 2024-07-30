//go:generate go run github.com/jmattheis/goverter/cmd/goverter gen .

package converter

import (
	"encoding/base64"
	"strconv"
	"time"

	"github.com/database-playground/backend/internal/models"
	"github.com/database-playground/backend/internal/services/gateway/openapi"
)

// goverter:converter
// goverter:matchIgnoreCase
// goverter:useZeroValueOnPointerInconsistency
// goverter:extend Int64ToString
// goverter:extend PInt64ToPString
// goverter:extend TimeToTime
type Converter interface {
	SchemaFromModel(in *models.Schema) openapi.Schema
	// goverter:enum:unknown Empty
	// goverter:enum:map DifficultyUnspecified Empty
	// goverter:enum:map DifficultyEasy Easy
	// goverter:enum:map DifficultyMedium Medium
	// goverter:enum:map DifficultyHard Hard
	DifficultyFromModel(in models.Difficulty) openapi.QuestionDifficulty
	QuestionFromModel(in *models.Question) openapi.Question
	QuestionsFromModel(in []*models.Question) openapi.Questions
	QuestionSolutionFromModel(in *models.QuestionSolution) openapi.QuestionSolution
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

func TimeToTime(in time.Time) time.Time {
	return in
}

func StringToID(in string) (int64, error) {
	return strconv.ParseInt(in, 10, 64)
}

func EncodeChallengeID(in string) string {
	return base64.URLEncoding.EncodeToString([]byte(in))
}

func DecodeChallengeID(in string) (string, error) {
	b, err := base64.URLEncoding.DecodeString(in)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
