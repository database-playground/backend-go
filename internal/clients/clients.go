// Package clients provides the client to interact with the microservices.
package clients

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"

	"github.com/database-playground/backend/gen/dbrunner/v1/dbrunnerv1connect"
	"github.com/database-playground/backend/gen/questionmanager/v1/questionmanagerv1connect"
	"go.uber.org/fx"
)

// NewConnectHTTPClient creates a new HTTP client to connect to the service.
func NewConnectHTTPClient(baseURL string) (*http.Client, error) {
	httpClient := &http.Client{}

	// add TLS certificate to the client if CLIENT_TLS_CERT_FILE is set
	if certFile := os.Getenv("CLIENT_TLS_CERT_FILE"); certFile != "" {
		cert, err := tls.LoadX509KeyPair(certFile, os.Getenv("CLIENT_TLS_KEY_FILE"))
		if err != nil {
			return nil, fmt.Errorf("load X509 key pair: %w", err)
		}
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
			},
		}

		// if TLS_CA_CERT_FILE is set, we should add it to the root CA chain
		if caFile := os.Getenv("TLS_CA_CERT_FILE"); caFile != "" {
			caCert, err := os.ReadFile(caFile)
			if err != nil {
				return nil, fmt.Errorf("read CA certificate file: %w", err)
			}
			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM(caCert)
			transport.TLSClientConfig.RootCAs = caCertPool
		}

		httpClient.Transport = transport
	}

	return httpClient, nil
}

var DBRunnerClientFxModule = fx.Module("dbrunner-client", fx.Provide(NewDBRunnerClient))

func NewDBRunnerClient() (dbrunnerv1connect.DbRunnerServiceClient, error) {
	baseURL := os.Getenv("DB_RUNNER_SERVICE_URL")
	if baseURL == "" {
		return nil, fmt.Errorf("DB_RUNNER_SERVICE_URL is not set")
	}

	httpClient, err := NewConnectHTTPClient(baseURL)
	if err != nil {
		return nil, fmt.Errorf("create HTTP client: %w", err)
	}

	return dbrunnerv1connect.NewDbRunnerServiceClient(httpClient, baseURL), nil
}

var QuestionManagerClientFxModule = fx.Module("question-manager-client", fx.Provide(NewQuestionManagerClient))

func NewQuestionManagerClient() (questionmanagerv1connect.QuestionManagerServiceClient, error) {
	baseURL := os.Getenv("QUESTION_MANAGER_SERVICE_URL")
	if baseURL == "" {
		return nil, fmt.Errorf("QUESTION_MANAGER_SERVICE_URL is not set")
	}

	httpClient, err := NewConnectHTTPClient(baseURL)
	if err != nil {
		return nil, fmt.Errorf("create HTTP client: %w", err)
	}

	return questionmanagerv1connect.NewQuestionManagerServiceClient(httpClient, baseURL), nil
}
