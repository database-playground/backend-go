package dbrunnerservice

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	dbrunnerv1 "github.com/database-playground/backend/gen/dbrunner/v1"
	"github.com/database-playground/backend/internal/dbrunner"
	"modernc.org/sqlite"
)

func (s *DBRunnerService) RunQuery(ctx context.Context, request *connect.Request[dbrunnerv1.RunQueryRequest]) (*connect.Response[dbrunnerv1.RunQueryResponse], error) {
	input := dbrunner.Input{
		Init:  request.Msg.GetSchema(),
		Query: request.Msg.GetQuery(),
	}

	// normalize input so it is cachable
	normalizedInput, err := input.Normalize()
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	// check if the output is existed; if so, return it.
	if outputHash, err := s.cacheModule.GetOutputHash(ctx, normalizedInput.Hash()); err == nil && s.cacheModule.HasOutput(ctx, outputHash) {
		return &connect.Response[dbrunnerv1.RunQueryResponse]{
			Msg: &dbrunnerv1.RunQueryResponse{
				ResponseType: &dbrunnerv1.RunQueryResponse_Id{
					Id: outputHash,
				},
			},
		}, nil
	}

	output, err := dbrunner.RunQuery(ctx, normalizedInput)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return &connect.Response[dbrunnerv1.RunQueryResponse]{
				Msg: &dbrunnerv1.RunQueryResponse{
					ResponseType: &dbrunnerv1.RunQueryResponse_Error{
						Error: "query timeout (takes more than 1 second)",
					},
				},
			}, nil
		}

		if errors.As(err, new(*sqlite.Error)) {
			return &connect.Response[dbrunnerv1.RunQueryResponse]{
				Msg: &dbrunnerv1.RunQueryResponse{
					ResponseType: &dbrunnerv1.RunQueryResponse_Error{
						Error: err.Error(),
					},
				},
			}, nil
		}

		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// cache the output
	id, err := s.cacheModule.WriteToCache(ctx, normalizedInput, output)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &connect.Response[dbrunnerv1.RunQueryResponse]{
		Msg: &dbrunnerv1.RunQueryResponse{
			ResponseType: &dbrunnerv1.RunQueryResponse_Id{
				Id: id,
			},
		},
	}, nil
}
