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
// goverter:extend UUIDToString
type Converter interface {
	// goverter:ignore state sizeCache unknownFields
	// goverter:map ID Id
	SchemaToProto(in *Schema) *questionmanagerv1.Schema

	// goverter:ignore state sizeCache unknownFields
	// goverter:map ID Id
	// goverter:map InitialSQL InitialSql
	SchemaInitialSQLToProto(in *SchemaInitialSQL) *questionmanagerv1.SchemaInitialSQL

	// goverter:ignore state sizeCache unknownFields
	// goverter:map ID Id
	// goverter:map SchemaID SchemaId
	QuestionToProto(in *Question) *questionmanagerv1.Question

	QuestionsToProto(in []*Question) []*questionmanagerv1.Question

	// goverter:enum:unknown Difficulty_DIFFICULTY_UNSPECIFIED
	// goverter:enum:map DifficultyUnspecified Difficulty_DIFFICULTY_UNSPECIFIED
	// goverter:enum:map DifficultyEasy Difficulty_DIFFICULTY_EASY
	// goverter:enum:map DifficultyMedium Difficulty_DIFFICULTY_MEDIUM
	// goverter:enum:map DifficultyHard Difficulty_DIFFICULTY_HARD
	DifficultyToProto(in Difficulty) questionmanagerv1.Difficulty
}

func TimeToTimestamp(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}

func UUIDToString(id uuid.UUID) string {
	return id.String()
}
