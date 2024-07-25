package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"connectrpc.com/connect"
	dbrunnerv1 "github.com/database-playground/backend/gen/dbrunner/v1"
	"github.com/database-playground/backend/gen/dbrunner/v1/dbrunnerv1connect"
)

func main() {
	baseURL := flag.String("base-url", "http://localhost:3000", "base URL of the service")
	schema := flag.String("schema", "", "schema SQL")
	query := flag.String("query", "", "query SQL")
	flag.Parse()

	if *schema == "" || *query == "" {
		panic("schema and query are required")
	}

	httpClient := http.DefaultClient

	// add TLS certificate to the client if CLIENT_TLS_CERT_FILE is set
	if certFile := os.Getenv("CLIENT_TLS_CERT_FILE"); certFile != "" {
		cert, err := tls.LoadX509KeyPair(certFile, os.Getenv("CLIENT_TLS_KEY_FILE"))
		if err != nil {
			panic(err)
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
				panic(err)
			}
			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM(caCert)
			transport.TLSClientConfig.RootCAs = caCertPool
		}

		httpClient.Transport = transport

		// change baseURL to https if TLS is enabled
		*baseURL = strings.Replace(*baseURL, "http://", "https://", 1)
	}

	client := dbrunnerv1connect.NewDbRunnerServiceClient(httpClient, *baseURL)
	mainQueryResponse, err := client.RunQuery(context.Background(), connect.NewRequest(&dbrunnerv1.RunQueryRequest{
		Schema: *schema,
		Query:  *query,
	}))
	if err != nil {
		panic(err)
	}

	if errMessage := mainQueryResponse.Msg.GetError(); errMessage != "" {
		panic(errMessage)
	}

	fmt.Printf("\x1b[90mINPUT_HASH: %v\x1b[0m\n", mainQueryResponse.Msg.GetId())

	// retrieve all the rows
	stream, err := client.RetrieveQuery(context.Background(), connect.NewRequest(&dbrunnerv1.RetrieveQueryRequest{
		Id: mainQueryResponse.Msg.GetId(),
	}))
	if err != nil {
		panic(err)
	}

	fmt.Printf("\x1b[90mOUTPUT_HASH: %v\x1b[0m\n", stream.ResponseHeader().Get("output-hash"))

	for stream.Receive() {
		switch msg := stream.Msg().GetKind().(type) {
		case *dbrunnerv1.RetrieveQueryResponse_Header:
			header := msg.Header

			fmt.Println("\x1b[1m" + strings.Join(header.Header, "\t") + "\x1b[0m")
		case *dbrunnerv1.RetrieveQueryResponse_Row:
			row := msg.Row

			for i, cell := range row.Cells {
				if i > 0 {
					fmt.Print("\t")
				}
				if cell.Value == nil {
					fmt.Print("\x1b[3mNULL\x1b[0m")
				} else {
					fmt.Print(*cell.Value)
				}

				if i == len(row.Cells)-1 {
					fmt.Println()
				}
			}
		}
	}
	if stream.Err() != nil {
		panic(stream.Err())
	}
}
