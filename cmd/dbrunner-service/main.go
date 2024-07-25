package main

import (
	"connectrpc.com/connect"
	"github.com/database-playground/backend/gen/dbrunner/v1/dbrunnerv1connect"
	httpservermodule "github.com/database-playground/backend/internal/modules/httpserver"
	redismodule "github.com/database-playground/backend/internal/modules/redis"
	slogmodule "github.com/database-playground/backend/internal/modules/slog"
	dbrunnerservice "github.com/database-playground/backend/internal/services/dbrunner"
	"go.uber.org/fx"
)

func main() {
	fx.New(slogmodule.FxOptions, redismodule.FxModule, dbrunnerservice.FxModule, fx.Provide(func(s *dbrunnerservice.Service) httpservermodule.HTTPHandler {
		return httpservermodule.WrapHTTPHandler[dbrunnerv1connect.DbRunnerServiceHandler](dbrunnerv1connect.NewDbRunnerServiceHandler, s, connect.WithRequireConnectProtocolHeader())
	}), httpservermodule.FxModule).Run()
}
