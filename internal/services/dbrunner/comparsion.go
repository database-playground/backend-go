package dbrunnerservice

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	dbrunnerv1 "github.com/database-playground/backend/gen/dbrunner/v1"
)

func (s *Service) AreQueriesOutputSame(ctx context.Context, request *connect.Request[dbrunnerv1.AreQueriesOutputSameRequest]) (*connect.Response[dbrunnerv1.AreQueriesOutputSameResponse], error) {
	leftHash := request.Msg.GetLeftId()
	rightHash := request.Msg.GetRightId()

	if leftHash == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("left_id is required"))
	}
	if rightHash == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("right_id is required"))
	}

	leftOutputHash, err := s.cacheModule.GetOutputHash(ctx, leftHash)
	if errors.Is(err, ErrNotFound) {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("left_id expired – re-query again!"))
	}
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	rightOutputHash, err := s.cacheModule.GetOutputHash(ctx, rightHash)
	if errors.Is(err, ErrNotFound) {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("right_id expired – re-query again!"))
	}
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &connect.Response[dbrunnerv1.AreQueriesOutputSameResponse]{
		Msg: &dbrunnerv1.AreQueriesOutputSameResponse{
			Same: leftOutputHash == rightOutputHash,
		},
	}, nil
}
