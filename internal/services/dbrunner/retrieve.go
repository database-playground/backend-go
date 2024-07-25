package dbrunnerservice

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	dbrunnerv1 "github.com/database-playground/backend/gen/dbrunner/v1"
	"github.com/samber/lo"
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

	// Set output-hash in the header.
	stream.ResponseHeader().Add("output-hash", outputHash)

	// Send header as the first packet
	if err := stream.Send(&dbrunnerv1.RetrieveQueryResponse{
		Kind: &dbrunnerv1.RetrieveQueryResponse_Header{
			Header: &dbrunnerv1.HeaderRow{
				Header: output.Header,
			},
		},
	}); err != nil {
		return err
	}

	// Send data rows one by one
	for _, row := range output.Data {
		if err := stream.Send(&dbrunnerv1.RetrieveQueryResponse{
			Kind: &dbrunnerv1.RetrieveQueryResponse_Row{
				Row: &dbrunnerv1.DataRow{
					Cells: lo.Map(row, func(cell *string, _ int) *dbrunnerv1.Cell {
						return &dbrunnerv1.Cell{
							Value: cell,
						}
					}),
				},
			},
		}); err != nil {
			return err
		}
	}

	return nil
}
