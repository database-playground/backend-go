package httpservermodule

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log/slog"
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

var FxModule = fx.Module("generic-http-server", fx.Provide(createTLSCertificate), fx.Provide(createTLSCertPool), fx.Invoke(func(handler HTTPHandler, cert *tls.Certificate, certPool *x509.CertPool, lc fx.Lifecycle) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	listenedOn := fmt.Sprintf("0.0.0.0:%s", port)

	srv := &http.Server{
		Addr:    listenedOn,
		Handler: handler,
	}
	if cert != nil {
		srv.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{*cert},
		}

		if certPool != nil {
			srv.TLSConfig.ClientCAs = certPool
			srv.TLSConfig.ClientAuth = tls.RequireAndVerifyClientCert
		}
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

func createTLSCertificate(logger slog.Logger) (*tls.Certificate, error) {
	certFile := os.Getenv("TLS_CERT_FILE")
	keyFile := os.Getenv("TLS_KEY_FILE")

	if certFile == "" || keyFile == "" {
		logger.Warn("TLS_CERT_FILE or TLS_KEY_FILE is not set, skipping TLS setup")
		return nil, nil
	}

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load TLS certificate: %w", err)
	}

	return &cert, nil
}

func createTLSCertPool(logger slog.Logger) (*x509.CertPool, error) {
	caFile := os.Getenv("TLS_CA_CERT_FILE")

	if caFile == "" {
		logger.Warn("TLS_CA_CERT_FILE is not set, skipping mTLS setup")
		return nil, nil
	}

	caCert, err := os.ReadFile(caFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA certificate: %w", err)
	}

	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to append CA certificate")
	}

	return pool, nil
}
