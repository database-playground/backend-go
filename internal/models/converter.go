//go:generate go run github.com/jmattheis/goverter/cmd/goverter gen .

package models

import (
	"time"

	questionmanagerv1 "github.com/database-playground/backend/gen/questionmanager/v1"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// goverter:converter
// goverter:extend TimeToTimestamp
// goverter:extend TimestampToTime
// goverter:extend UUIDToString
// goverter:extend StringToUUID
type Converter interface {
	// goverter:ignore state sizeCache unknownFields
	// goverter:map ID Id
	SchemaToProto(in *Schema) *questionmanagerv1.Schema

	// goverter:map Id ID
	SchemaFromProto(in *questionmanagerv1.Schema) *Schema

	// goverter:ignore state sizeCache unknownFields
	// goverter:map ID Id
	// goverter:map InitialSQL InitialSql
	SchemaInitialSQLToProto(in *SchemaInitialSQL) *questionmanagerv1.SchemaInitialSQL

	// goverter:map Id ID
	// goverter:map InitialSql InitialSQL
	SchemaInitialSQLFromProto(in *questionmanagerv1.SchemaInitialSQL) *SchemaInitialSQL

	// goverter:ignore state sizeCache unknownFields
	// goverter:map ID Id
	// goverter:map SchemaID SchemaId
	QuestionToProto(in *Question) *questionmanagerv1.Question

	// goverter:map Id ID
	// goverter:map SchemaId SchemaID
	QuestionFromProto(in *questionmanagerv1.Question) *Question

	QuestionsToProto(in []*Question) []*questionmanagerv1.Question

	QuestionsFromProto(in []*questionmanagerv1.Question) []*Question

	// goverter:enum:unknown Difficulty_DIFFICULTY_UNSPECIFIED
	// goverter:enum:map DifficultyUnspecified Difficulty_DIFFICULTY_UNSPECIFIED
	// goverter:enum:map DifficultyEasy Difficulty_DIFFICULTY_EASY
	// goverter:enum:map DifficultyMedium Difficulty_DIFFICULTY_MEDIUM
	// goverter:enum:map DifficultyHard Difficulty_DIFFICULTY_HARD
	DifficultyToProto(in Difficulty) questionmanagerv1.Difficulty

	// goverter:enum:unknown DifficultyUnspecified
	// goverter:enum:map Difficulty_DIFFICULTY_UNSPECIFIED DifficultyUnspecified
	// goverter:enum:map Difficulty_DIFFICULTY_EASY DifficultyEasy
	// goverter:enum:map Difficulty_DIFFICULTY_MEDIUM DifficultyMedium
	// goverter:enum:map Difficulty_DIFFICULTY_HARD DifficultyHard
	DifficultyFromProto(in questionmanagerv1.Difficulty) Difficulty

	// goverter:ignore state sizeCache unknownFields
	// goverter:map ID Id
	QuestionAnswerToProto(in *QuestionAnswer) *questionmanagerv1.QuestionAnswer

	// goverter:map Id ID
	QuestionAnswerFromProto(in *questionmanagerv1.QuestionAnswer) *QuestionAnswer

	// goverter:ignore state sizeCache unknownFields
	// goverter:map ID Id
	QuestionSolutionToProto(in *QuestionSolution) *questionmanagerv1.QuestionSolution

	// goverter:map Id ID
	QuestionSolutionFromProto(in *questionmanagerv1.QuestionSolution) *QuestionSolution
}

func TimeToTimestamp(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}

func TimestampToTime(t *timestamppb.Timestamp) time.Time {
	return t.AsTime()
}

func UUIDToString(id uuid.UUID) string {
	return id.String()
}

func StringToUUID(id string) uuid.UUID {
	uuidValue, _ := uuid.Parse(id)
	return uuidValue
}
