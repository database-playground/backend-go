package dbrunnerservice

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	commonv1 "github.com/database-playground/backend/gen/common/v1"
	dbrunnerv1 "github.com/database-playground/backend/gen/dbrunner/v1"
)

func (s *Service) RetrieveQuery(ctx context.Context, request *connect.Request[dbrunnerv1.RetrieveQueryRequest], stream *connect.ServerStream[dbrunnerv1.RetrieveQueryResponse]) error {
	if request.Msg.GetId() == "" {
		return connect.NewError(connect.CodeInvalidArgument, errors.New("id is required"))
	}

	outputHash, err := s.cacheModule.GetOutputHash(ctx, request.Msg.GetId())
	if errors.Is(err, ErrNotFound) {
		return connect.NewError(connect.CodeNotFound, errors.New("id expired – re-query again!"))
	}
	if err != nil {
		return connect.NewError(connect.CodeInternal, err)
	}

	output, err := s.cacheModule.GetOutput(ctx, outputHash)
	if errors.Is(err, ErrNotFound) {
		return connect.NewError(connect.CodeNotFound, errors.New("output expired – re-query again!"))
	}
	if err != nil {
		return connect.NewError(connect.CodeInternal, err)
	}

	for _, resultRow := range output.Result {
		rpcRow := make([]*commonv1.OptionalStringPair, 0, len(resultRow))

		for _, value := range resultRow {
			rpcRow = append(rpcRow, &commonv1.OptionalStringPair{
				Key:   value.Column,
				Value: value.Value,
			})
		}

		if err := stream.Send(&dbrunnerv1.RetrieveQueryResponse{
			Row: rpcRow,
		}); err != nil {
			return err
		}
	}

	return nil
}
