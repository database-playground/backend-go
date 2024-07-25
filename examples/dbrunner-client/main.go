package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
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

	client := dbrunnerv1connect.NewDbRunnerServiceClient(http.DefaultClient, *baseURL)
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
