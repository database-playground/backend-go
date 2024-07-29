package gatewayservice

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/database-playground/backend/gen/dbrunner/v1/dbrunnerv1connect"
	"github.com/database-playground/backend/gen/questionmanager/v1/questionmanagerv1connect"
	"github.com/database-playground/backend/internal/models"
	pbgenerated "github.com/database-playground/backend/internal/models/generated"
	"go.uber.org/fx"

	"github.com/database-playground/backend/internal/services/gateway/converter"
	modelgenerated "github.com/database-playground/backend/internal/services/gateway/converter/generated"
	"github.com/database-playground/backend/internal/services/gateway/openapi"
)

//go:embed openapi/openapi.yaml
var openapiSpec []byte

var FxModule = fx.Module("gateway-service", fx.Provide(NewServer), fx.Invoke(func(server openapi.StrictServerInterface, lc fx.Lifecycle) {
	mux := http.NewServeMux()

	// serve openapi spec
	mux.HandleFunc("GET /openapi.yaml", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(openapiSpec)
	})

	handler := openapi.HandlerFromMux(openapi.NewStrictHandler(server, []openapi.StrictMiddlewareFunc{}), mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	certFile := os.Getenv("PUBLIC_TLS_CERT_FILE")
	keyFile := os.Getenv("PUBLIC_TLS_KEY_FILE")

	s := &http.Server{
		Handler: handler,
		Addr:    "0.0.0.0:" + port,
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				// no mTLS is needed in gateway service
				if certFile != "" && keyFile != "" {
					_ = s.ListenAndServeTLS(certFile, keyFile)
				} else {
					_ = s.ListenAndServe()
				}
			}()
			fmt.Printf("starting server at %s\n", s.Addr)
			return nil
		},
		OnStop: func(context.Context) error {
			return s.Shutdown(context.Background())
		},
	})
}))

type ServerParam struct {
	fx.In

	Logger                *slog.Logger
	QuestionManagerClient questionmanagerv1connect.QuestionManagerServiceClient
	DBRunnerClient        dbrunnerv1connect.DbRunnerServiceClient
}

type Server struct {
	logger *slog.Logger

	questionManagerService questionmanagerv1connect.QuestionManagerServiceClient
	dbrunnerService        dbrunnerv1connect.DbRunnerServiceClient

	pbConverter    models.Converter
	modelConverter converter.Converter
}

func NewServer(param ServerParam) openapi.StrictServerInterface {
	return &Server{
		logger: param.Logger,

		questionManagerService: param.QuestionManagerClient,
		dbrunnerService:        param.DBRunnerClient,

		pbConverter:    &pbgenerated.ConverterImpl{},
		modelConverter: &modelgenerated.ConverterImpl{},
	}
}
