package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	dbrunnerv1 "github.com/database-playground/backend/gen/dbrunner/v1"
	"github.com/database-playground/backend/gen/dbrunner/v1/dbrunnerv1connect"
)

func main() {
	baseURL := flag.String("base-url", "http://localhost:3000", "base URL of the service")
	schema := flag.String("schema", "", "schema SQL")
	query := flag.String("query", "", "query SQL")
	query2 := flag.String("query2", "", "query SQL to compare")
	flag.Parse()

	if *schema == "" || *query == "" || *query2 == "" {
		panic("schema, query and query2 are required")
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

	secondaryQueryResponse, err := client.RunQuery(context.Background(), connect.NewRequest(&dbrunnerv1.RunQueryRequest{
		Schema: *schema,
		Query:  *query2,
	}))
	if err != nil {
		panic(err)
	}
	if errMessage := secondaryQueryResponse.Msg.GetError(); errMessage != "" {
		panic(errMessage)
	}

	fmt.Printf("\x1b[90mINPUT_1_HASH: %s\x1b[0m\n", mainQueryResponse.Msg.GetId())
	fmt.Printf("\x1b[90mINPUT_2_HASH: %s\x1b[0m\n", secondaryQueryResponse.Msg.GetId())

	sameQueryResponse, err := client.AreQueriesOutputSame(context.Background(), connect.NewRequest(&dbrunnerv1.AreQueriesOutputSameRequest{
		LeftId:  mainQueryResponse.Msg.GetId(),
		RightId: secondaryQueryResponse.Msg.GetId(),
	}))
	if err != nil {
		panic(err)
	}

	fmt.Printf("\x1b[90mOUTPUT_SAME? %v\x1b[0m\n", sameQueryResponse.Msg.GetSame())
}
