package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/database-playground/backend/gen/dbrunner/v1/dbrunnerv1connect"
	redismodule "github.com/database-playground/backend/internal/modules/redis"
	slogmodule "github.com/database-playground/backend/internal/modules/slog"
	dbrunnerservice "github.com/database-playground/backend/internal/services/dbrunner"
	"go.uber.org/fx"
)

func main() {
	fx.New(slogmodule.FxOptions, redismodule.FxModule, dbrunnerservice.FxModule, fx.Invoke(func(service *dbrunnerservice.Service, lc fx.Lifecycle) {
		port := os.Getenv("PORT")
		if port == "" {
			port = "3000"
		}
		listenedOn := fmt.Sprintf("0.0.0.0:%s", port)

		path, handler := dbrunnerv1connect.NewDbRunnerServiceHandler(service)
		srv := &http.Server{
			Addr:    listenedOn,
			Handler: handler,
		}

		lc.Append(fx.Hook{
			OnStart: func(context.Context) error {
				ln, err := net.Listen("tcp", srv.Addr)
				if err != nil {
					return err
				}
				go func() {
					_ = srv.Serve(ln)
				}()
				fmt.Printf("starting server at %s%s\n", listenedOn, path)
				return nil
			},
			OnStop: func(ctx context.Context) error {
				fmt.Println("stopping server")
				go func() {
					_ = srv.Shutdown(ctx)
				}()
				return nil
			},
		})
	})).Run()
}
