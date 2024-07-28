package questionmanagerservice

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	questionmanagerv1 "github.com/database-playground/backend/gen/questionmanager/v1"
	"github.com/database-playground/backend/internal/database"
)

func (s *Service) ListQuestions(ctx context.Context, request *connect.Request[questionmanagerv1.ListQuestionsRequest]) (*connect.Response[questionmanagerv1.ListQuestionsResponse], error) {
	questions, err := s.db.ListQuestions(ctx, database.ListQuestionsParams{
		Cursor: database.CursorFromProto(request.Msg.Cursor),
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	questionsPb := s.converter.QuestionsToProto(questions)

	return &connect.Response[questionmanagerv1.ListQuestionsResponse]{
		Msg: &questionmanagerv1.ListQuestionsResponse{
			Questions: questionsPb,
		},
	}, nil
}

func (s *Service) GetQuestion(ctx context.Context, request *connect.Request[questionmanagerv1.GetQuestionRequest]) (*connect.Response[questionmanagerv1.GetQuestionResponse], error) {
	question, err := s.db.GetQuestion(ctx, request.Msg.GetId())
	if errors.Is(err, database.ErrNotFound) {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	questionPb := s.converter.QuestionToProto(question)

	return &connect.Response[questionmanagerv1.GetQuestionResponse]{
		Msg: &questionmanagerv1.GetQuestionResponse{
			Question: questionPb,
		},
	}, nil
}
