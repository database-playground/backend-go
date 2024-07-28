package questionmanagerservice

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	questionmanagerv1 "github.com/database-playground/backend/gen/questionmanager/v1"
	"github.com/database-playground/backend/internal/database"
)

func (s *Service) GetSchema(ctx context.Context, request *connect.Request[questionmanagerv1.GetSchemaRequest]) (*connect.Response[questionmanagerv1.GetSchemaResponse], error) {
	schema, err := s.db.GetSchema(ctx, request.Msg.GetId())
	if errors.Is(err, database.ErrNotFound) {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	schemaPb := s.converter.SchemaToProto(schema)
	return &connect.Response[questionmanagerv1.GetSchemaResponse]{
		Msg: &questionmanagerv1.GetSchemaResponse{
			Schema: schemaPb,
		},
	}, nil
}

func (s *Service) GetSchemaInitialSQL(ctx context.Context, request *connect.Request[questionmanagerv1.GetSchemaInitialSQLRequest]) (*connect.Response[questionmanagerv1.GetSchemaInitialSQLResponse], error) {
	schemaInitialSQL, err := s.db.GetSchemaInitialSQL(ctx, request.Msg.GetId())
	if errors.Is(err, database.ErrNotFound) {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	schemaInitialSQLPb := s.converter.SchemaInitialSQLToProto(schemaInitialSQL)
	return &connect.Response[questionmanagerv1.GetSchemaInitialSQLResponse]{
		Msg: &questionmanagerv1.GetSchemaInitialSQLResponse{
			SchemaInitialSql: schemaInitialSQLPb,
		},
	}, nil
}
