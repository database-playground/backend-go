//go:generate go run github.com/jmattheis/goverter/cmd/goverter gen .

package models

import (
	"time"

	questionmanagerv1 "github.com/database-playground/backend/gen/questionmanager/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// goverter:converter
// goverter:extend TimeToTimestamp
type Converter interface {
	// goverter:ignore state sizeCache unknownFields
	// goverter:map ID Id
	SchemaToProto(in *Schema) *questionmanagerv1.Schema

	// goverter:ignore state sizeCache unknownFields
	// goverter:map ID Id
	// goverter:map InitialSQL InitialSql
	SchemaInitialSQLToProto(in *SchemaInitialSQL) *questionmanagerv1.SchemaInitialSQL
}

func TimeToTimestamp(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}
