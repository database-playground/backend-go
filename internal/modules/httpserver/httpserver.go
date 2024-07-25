package httpservermodule

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"

	"connectrpc.com/connect"
	"go.uber.org/fx"
)

// HTTPHandler is a struct that contains the path and handler of an HTTP handler.
//
// It can be wrapped from NewServiceHandler from connect package.
type HTTPHandler struct {
	http.Handler
	RpcPath string
}

// WrapHTTPHandler is a function that returns an HTTPHandler struct.
//
// It accepts functions like [dbrunnerv1connect.NewDbRunnerServiceHandler]
// and wraps it to return an [HTTPHandler] struct.
func WrapHTTPHandler[S any](f func(S, ...connect.HandlerOption) (string, http.Handler), service S, handlerOptions ...connect.HandlerOption) HTTPHandler {
	path, handler := f(service)
	return HTTPHandler{
		RpcPath: path,
		Handler: handler,
	}
}

var FxModule = fx.Module("generic-http-server", fx.Invoke(func(handler HTTPHandler, lc fx.Lifecycle) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	listenedOn := fmt.Sprintf("0.0.0.0:%s", port)

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
			fmt.Printf("starting server at %s%s\n", listenedOn, handler.RpcPath)
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
}))
