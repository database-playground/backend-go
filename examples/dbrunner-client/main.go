package main

import (
	"context"
	"flag"
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
	stream, err := client.RunQuery(context.Background(), connect.NewRequest(&dbrunnerv1.RunQueryRequest{
		Schema: *schema,
		Query:  *query,
	}))
	if err != nil {
		panic(err)
	}

	headerPrinted := false

	for stream.Receive() {
		resp := stream.Msg()

		header := []string{}
		content := []string{}
		for _, row := range resp.Rows {
			header = append(header, row.Key)
			if row.Value != nil {
				content = append(content, *row.Value)
			} else {
				content = append(content, "\x1b[3mNULL\x1b[0m")
			}
		}

		if !headerPrinted {
			headerPrinted = true
			println("\x1b[1m" + strings.Join(header, "\t") + "\x1b[0m")
		}
		println(strings.Join(content, "\t"))
	}
	if stream.Err() != nil {
		panic(stream.Err())
	}
}
