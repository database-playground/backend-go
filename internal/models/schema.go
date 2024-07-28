package models

import (
	"time"
)

// Schema represents a database schema that can be applied to a question.
type Schema struct {
	ID string `json:"id" db:"schema_id"`
	// Picture is a URL to a picture of the schema relationships.
	Picture *string `json:"picture,omitempty"`
	// Description is a description of the schema.
	Description string `json:"description,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SchemaInitialSQL struct {
	ID         string `json:"id" db:"schema_id"`
	InitialSQL string `json:"inital_sql" db:"initial_sql"`
}

// Difficulty represents the difficulty of a question.
type Difficulty string

const (
	DifficultyUnspecified Difficulty = ""
	DifficultyEasy        Difficulty = "easy"
	DifficultyMedium      Difficulty = "medium"
	DifficultyHard        Difficulty = "hard"
)

type Question struct {
	ID       int64  `json:"id" db:"question_id"`
	SchemaID string `json:"schema_id"`

	Type       string     `json:"type"`
	Difficulty Difficulty `json:"difficulty"`

	Title       string `json:"title"`
	Description string `json:"description"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type QuestionAnswer struct {
	ID int64 `json:"id"`

	// Answer is the correct answer to the question.
	Answer string `json:"answer"`
}

type QuestionSolution struct {
	ID int64 `json:"id"`

	// SolutionVideo is a URL to a video that explains the solution.
	SolutionVideo *string `json:"solution_video,omitempty"`
}
