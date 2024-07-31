package gatewayservice

import (
	"context"
	"log/slog"

	"connectrpc.com/connect"
	commonv1 "github.com/database-playground/backend/gen/common/v1"
	dbrunnerv1 "github.com/database-playground/backend/gen/dbrunner/v1"
	questionmanagerv1 "github.com/database-playground/backend/gen/questionmanager/v1"
	"github.com/database-playground/backend/internal/services/gateway/converter"
	"github.com/database-playground/backend/internal/services/gateway/openapi"
)

var _ openapi.StrictServerInterface = (*Server)(nil)

// #region Health Check

// GetHealthz implements openapi.StrictServerInterface.
func (s *Server) GetHealthz(context.Context, openapi.GetHealthzRequestObject) (openapi.GetHealthzResponseObject, error) {
	return openapi.GetHealthz200Response{}, nil
}

// #region Questions

// GetQuestions implements StrictServerInterface.
func (s *Server) GetQuestions(ctx context.Context, request openapi.GetQuestionsRequestObject) (openapi.GetQuestionsResponseObject, error) {
	response, err := s.questionManagerService.ListQuestions(ctx, &connect.Request[questionmanagerv1.ListQuestionsRequest]{
		Msg: &questionmanagerv1.ListQuestionsRequest{
			Cursor: &commonv1.Cursor{
				Limit:  request.Params.Limit,
				Offset: request.Params.Offset,
			},
		},
	})
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to fetch questions", slog.Any("error", err), slog.Any("request", request))
		return openapi.GetQuestions500JSONResponse{
			ErrorJSONResponse: openapi.ErrorJSONResponse{
				Message: "Failed to fetch questions.",
			},
		}, nil
	}

	questionsModel := s.pbConverter.QuestionsFromProto(response.Msg.GetQuestions())
	questionsResponse := s.modelConverter.QuestionsFromModel(questionsModel)

	return openapi.GetQuestions200JSONResponse(questionsResponse), nil
}

// GetQuestionsId implements StrictServerInterface.
func (s *Server) GetQuestionsId(ctx context.Context, request openapi.GetQuestionsIdRequestObject) (openapi.GetQuestionsIdResponseObject, error) {
	id, err := converter.StringToID(request.Id)
	if err != nil {
		return openapi.GetQuestionsId400JSONResponse{
			BadRequestErrorJSONResponse: openapi.BadRequestErrorJSONResponse{
				Message: "Invalid ID.",
			},
		}, nil
	}

	response, err := s.questionManagerService.GetQuestion(ctx, &connect.Request[questionmanagerv1.GetQuestionRequest]{
		Msg: &questionmanagerv1.GetQuestionRequest{
			Id: id,
		},
	})
	if connect.CodeOf(err) == connect.CodeNotFound {
		return openapi.GetQuestionsId404JSONResponse{
			NoSuchResourceErrorJSONResponse: openapi.NoSuchResourceErrorJSONResponse{
				Message: "Question not found.",
			},
		}, nil
	}
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to fetch question", slog.Any("error", err), slog.Any("request", request))
		return openapi.GetQuestionsId500JSONResponse{
			ErrorJSONResponse: openapi.ErrorJSONResponse{
				Message: "Failed to fetch question.",
			},
		}, nil
	}

	questionModel := s.pbConverter.QuestionFromProto(response.Msg.GetQuestion())
	questionResponse := s.modelConverter.QuestionFromModel(questionModel)

	return openapi.GetQuestionsId200JSONResponse(questionResponse), nil
}

// GetQuestionsIdSolution implements StrictServerInterface.
func (s *Server) GetQuestionsIdSolution(ctx context.Context, request openapi.GetQuestionsIdSolutionRequestObject) (openapi.GetQuestionsIdSolutionResponseObject, error) {
	id, err := converter.StringToID(request.Id)
	if err != nil {
		return openapi.GetQuestionsIdSolution400JSONResponse{
			BadRequestErrorJSONResponse: openapi.BadRequestErrorJSONResponse{
				Message: "Invalid ID.",
			},
		}, nil
	}

	response, err := s.questionManagerService.GetQuestionSolution(ctx, &connect.Request[questionmanagerv1.GetQuestionSolutionRequest]{
		Msg: &questionmanagerv1.GetQuestionSolutionRequest{
			Id: id,
		},
	})
	if connect.CodeOf(err) == connect.CodeNotFound {
		return openapi.GetQuestionsIdSolution404JSONResponse{
			NoSuchResourceErrorJSONResponse: openapi.NoSuchResourceErrorJSONResponse{
				Message: "Solution not found.",
			},
		}, nil
	}
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to fetch solution", slog.Any("error", err), slog.Any("request", request))
		return openapi.GetQuestionsIdSolution500JSONResponse{
			ErrorJSONResponse: openapi.ErrorJSONResponse{
				Message: "Failed to fetch solution.",
			},
		}, nil
	}

	solutionModel := s.pbConverter.QuestionSolutionFromProto(response.Msg.GetQuestionSolution())
	solutionResponse := s.modelConverter.QuestionSolutionFromModel(solutionModel)

	return openapi.GetQuestionsIdSolution200JSONResponse(solutionResponse), nil
}

// #region Question Challenge

// GetChallenge implements openapi.StrictServerInterface.
func (s *Server) GetChallengesId(ctx context.Context, request openapi.GetChallengesIdRequestObject) (openapi.GetChallengesIdResponseObject, error) {
	tc, err := converter.DecodeChallengeID(request.Id)
	if err != nil || tc == nil {
		return openapi.GetChallengesId400JSONResponse{
			BadRequestErrorJSONResponse: openapi.BadRequestErrorJSONResponse{
				Message: "Invalid challenge ID.",
			},
		}, nil
	}

	response, err := s.dbrunnerService.RetrieveQuery(ctx, &connect.Request[dbrunnerv1.RetrieveQueryRequest]{
		Msg: &dbrunnerv1.RetrieveQueryRequest{
			Id: tc.ChallengeID,
		},
	})
	if connect.CodeOf(err) == connect.CodeNotFound {
		return openapi.GetChallengesId404JSONResponse{
			NoSuchResourceErrorJSONResponse: openapi.NoSuchResourceErrorJSONResponse{
				Message: "Challenge not found or is expired.",
			},
		}, nil
	}
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to fetch challenge", slog.Any("error", err), slog.Any("request", request))
		return openapi.GetChallengesId500JSONResponse{
			ErrorJSONResponse: openapi.ErrorJSONResponse{
				Message: "Failed to fetch challenge.",
			},
		}, nil
	}

	var header []string
	var rows [][]*string

	for response.Receive() {
		switch messageKind := response.Msg().Kind.(type) {
		case *dbrunnerv1.RetrieveQueryResponse_Header:
			header = messageKind.Header.GetHeader()
		case *dbrunnerv1.RetrieveQueryResponse_Row:
			var row []*string
			cell := messageKind.Row.GetCells()
			for _, cell := range cell {
				row = append(row, cell.Value)
			}
			rows = append(rows, row)
		}
	}
	if response.Err() != nil {
		if connect.CodeOf(response.Err()) == connect.CodeNotFound {
			return openapi.GetChallengesId404JSONResponse{
				NoSuchResourceErrorJSONResponse: openapi.NoSuchResourceErrorJSONResponse{
					Message: "Challenge not found or is expired.",
				},
			}, nil
		}

		s.logger.ErrorContext(ctx, "Failed to fetch challenge", slog.Any("error", response.Err()), slog.Any("request", request))
		return openapi.GetChallengesId500JSONResponse{
			ErrorJSONResponse: openapi.ErrorJSONResponse{
				Message: "Failed to fetch challenge.",
			},
		}, nil
	}

	return openapi.GetChallengesId200JSONResponse{
		Header: header,
		Rows:   rows,
	}, nil
}

// PostChallenge implements openapi.StrictServerInterface.
func (s *Server) PostChallenges(ctx context.Context, request openapi.PostChallengesRequestObject) (openapi.PostChallengesResponseObject, error) {
	questionID, err := converter.StringToID(request.Body.QuestionID)
	if err != nil {
		return openapi.PostChallenges400JSONResponse{
			BadRequestErrorJSONResponse: openapi.BadRequestErrorJSONResponse{
				Message: "Invalid question ID.",
			},
		}, nil
	}

	questionResponse, err := s.questionManagerService.GetQuestion(ctx, &connect.Request[questionmanagerv1.GetQuestionRequest]{
		Msg: &questionmanagerv1.GetQuestionRequest{
			Id: questionID,
		},
	})
	if connect.CodeOf(err) == connect.CodeNotFound {
		return openapi.PostChallenges404JSONResponse{
			NoSuchResourceErrorJSONResponse: openapi.NoSuchResourceErrorJSONResponse{
				Message: "Question not found.",
			},
		}, nil
	}
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to fetch question", slog.Any("error", err), slog.Any("request", request))
		return openapi.PostChallenges500JSONResponse{
			ErrorJSONResponse: openapi.ErrorJSONResponse{
				Message: "Failed to fetch question.",
			},
		}, nil
	}

	schemaInitialSQLResponse, err := s.questionManagerService.GetSchemaInitialSQL(ctx, &connect.Request[questionmanagerv1.GetSchemaInitialSQLRequest]{
		Msg: &questionmanagerv1.GetSchemaInitialSQLRequest{
			Id: questionResponse.Msg.GetQuestion().GetSchemaId(),
		},
	})
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to fetch initial SQL", slog.Any("error", err), slog.Any("request", request))
		return openapi.PostChallenges500JSONResponse{
			ErrorJSONResponse: openapi.ErrorJSONResponse{
				Message: "Failed to fetch initial SQL.",
			},
		}, nil
	}

	// execute question
	queryResponse, err := s.dbrunnerService.RunQuery(ctx, &connect.Request[dbrunnerv1.RunQueryRequest]{
		Msg: &dbrunnerv1.RunQueryRequest{
			Schema: schemaInitialSQLResponse.Msg.GetSchemaInitialSql().GetInitialSql(),
			Query:  request.Body.Query,
		},
	})
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to execute query", slog.Any("error", err), slog.Any("request", request))
		return openapi.PostChallenges500JSONResponse{
			ErrorJSONResponse: openapi.ErrorJSONResponse{
				Message: "Failed to execute query (not user-side error).",
			},
		}, nil
	}
	if queryResponse.Msg.GetError() != "" {
		return openapi.PostChallenges422JSONResponse{
			UnprocessableEntityErrorJSONResponse: openapi.UnprocessableEntityErrorJSONResponse{
				Message: queryResponse.Msg.GetError(),
			},
		}, nil
	}

	// hash challenge ID so we can push it to URL
	base64ChallengeID := converter.EncodeChallengeID(converter.TransferableChallengeID{
		QuestionID:  questionID,
		ChallengeID: queryResponse.Msg.GetId(),
	})

	return openapi.PostChallenges200JSONResponse{
		ChallengeID: base64ChallengeID,
	}, nil
}

// GetChallengesIdCompare implements openapi.StrictServerInterface.
func (s *Server) GetChallengesIdCompare(ctx context.Context, request openapi.GetChallengesIdCompareRequestObject) (openapi.GetChallengesIdCompareResponseObject, error) {
	tc, err := converter.DecodeChallengeID(request.Id)
	if err != nil || tc == nil {
		return openapi.GetChallengesIdCompare400JSONResponse{
			BadRequestErrorJSONResponse: openapi.BadRequestErrorJSONResponse{
				Message: "Invalid challenge ID.",
			},
		}, nil
	}

	answer, err := s.questionManagerService.GetQuestionAnswer(ctx, &connect.Request[questionmanagerv1.GetQuestionAnswerRequest]{
		Msg: &questionmanagerv1.GetQuestionAnswerRequest{
			Id: tc.QuestionID,
		},
	})
	if connect.CodeOf(err) == connect.CodeNotFound {
		return openapi.GetChallengesIdCompare404JSONResponse{
			NoSuchResourceErrorJSONResponse: openapi.NoSuchResourceErrorJSONResponse{
				Message: "Answer not found.",
			},
		}, nil
	}
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to fetch answer", slog.Any("error", err), slog.Any("request", request))
		return openapi.GetChallengesIdCompare500JSONResponse{
			ErrorJSONResponse: openapi.ErrorJSONResponse{
				Message: "Failed to fetch answer.",
			},
		}, nil
	}

	answerResponse, err := s.dbrunnerService.RunQuery(ctx, &connect.Request[dbrunnerv1.RunQueryRequest]{
		Msg: &dbrunnerv1.RunQueryRequest{
			Schema: answer.Msg.QuestionAnswer.GetSchema(),
			Query:  answer.Msg.QuestionAnswer.GetAnswer(),
		},
	})
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to execute answer", slog.Any("error", err), slog.Any("request", request))
		return openapi.GetChallengesIdCompare500JSONResponse{
			ErrorJSONResponse: openapi.ErrorJSONResponse{
				Message: "Failed to execute answer. The answer is incorrect.",
			},
		}, nil
	}
	if answerResponse.Msg.GetError() != "" {
		s.logger.ErrorContext(ctx, "Failed to execute answer", slog.Any("error", answerResponse.Msg.GetError()), slog.Any("request", request))
		return openapi.GetChallengesIdCompare500JSONResponse{
			ErrorJSONResponse: openapi.ErrorJSONResponse{
				Message: "Failed to execute answer. The answer is incorrect.",
			},
		}, nil
	}

	sameResponse, err := s.dbrunnerService.AreQueriesOutputSame(ctx, &connect.Request[dbrunnerv1.AreQueriesOutputSameRequest]{
		Msg: &dbrunnerv1.AreQueriesOutputSameRequest{
			LeftId:  answerResponse.Msg.GetId(),
			RightId: tc.ChallengeID,
		},
	})
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to compare answers", slog.Any("error", err), slog.Any("request", request))
		return openapi.GetChallengesIdCompare500JSONResponse{
			ErrorJSONResponse: openapi.ErrorJSONResponse{
				Message: "Failed to compare answers.",
			},
		}, nil
	}

	return openapi.GetChallengesIdCompare200JSONResponse{
		Same: sameResponse.Msg.GetSame(),
	}, nil
}

// #region Schema

// GetSchemasId implements StrictServerInterface.
func (s *Server) GetSchemasId(ctx context.Context, request openapi.GetSchemasIdRequestObject) (openapi.GetSchemasIdResponseObject, error) {
	response, err := s.questionManagerService.GetSchema(ctx, &connect.Request[questionmanagerv1.GetSchemaRequest]{
		Msg: &questionmanagerv1.GetSchemaRequest{
			Id: request.Id,
		},
	})
	if connect.CodeOf(err) == connect.CodeNotFound {
		return openapi.GetSchemasId404JSONResponse{
			NoSuchResourceErrorJSONResponse: openapi.NoSuchResourceErrorJSONResponse{
				Message: "Schema not found.",
			},
		}, nil
	}
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to fetch schema", slog.Any("error", err), slog.Any("request", request))
		return openapi.GetSchemasId500JSONResponse{
			ErrorJSONResponse: openapi.ErrorJSONResponse{
				Message: "Failed to fetch schema.",
			},
		}, nil
	}

	schemaModel := s.pbConverter.SchemaFromProto(response.Msg.GetSchema())
	schemaResponse := s.modelConverter.SchemaFromModel(schemaModel)

	return openapi.GetSchemasId200JSONResponse(schemaResponse), nil
}
