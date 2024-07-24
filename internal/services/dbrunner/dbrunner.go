package dbrunnerservice

import (
	"context"
	"strings"

	"connectrpc.com/connect"
	commonv1 "github.com/database-playground/backend/gen/common/v1"
	dbrunnerv1 "github.com/database-playground/backend/gen/dbrunner/v1"
	"github.com/database-playground/backend/gen/dbrunner/v1/dbrunnerv1connect"
	"github.com/database-playground/backend/internal/dbrunner"
	"github.com/samber/lo"
)

type DBRunnerService struct {
	dbrunnerv1connect.UnimplementedDbRunnerServiceHandler
}

func NewDBRunnerService() *DBRunnerService {
	return &DBRunnerService{}
}

func (s *DBRunnerService) RunQuery(ctx context.Context, request *connect.Request[dbrunnerv1.RunQueryRequest], stream *connect.ServerStream[dbrunnerv1.RunQueryResponse]) error {
	input := dbrunner.Input{
		Init:  request.Msg.GetSchema(),
		Query: request.Msg.GetQuery(),
	}
	output, err := dbrunner.RunQuery(ctx, input)
	if err != nil {
		if strings.Contains(err.Error(), "query timeout") {
			return connect.NewError(connect.CodeAborted, err)
		}

		if strings.Contains(err.Error(), "exec init") ||
			strings.Contains(err.Error(), "query") ||
			strings.Contains(err.Error(), "sqlite error") {
			return connect.NewError(connect.CodeInvalidArgument, err)
		}

		return connect.NewError(connect.CodeInternal, err)
	}

	for _, result := range output.Result {
		err := stream.Send(&dbrunnerv1.RunQueryResponse{
			Rows: lo.Map(result, func(r struct {
				Column string
				Value  *string
			}, _ int,
			) *commonv1.OptionalStringPair {
				return &commonv1.OptionalStringPair{
					Key:   r.Column,
					Value: r.Value,
				}
			}),
		})
		if err != nil {
			return connect.NewError(connect.CodeDataLoss, err)
		}
	}

	return nil
}
