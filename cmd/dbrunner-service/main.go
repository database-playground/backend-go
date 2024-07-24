package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/database-playground/backend/gen/dbrunner/v1/dbrunnerv1connect"
	dbrunnerservice "github.com/database-playground/backend/internal/services/dbrunner"
)

func main() {
	path, handler := dbrunnerv1connect.NewDbRunnerServiceHandler(dbrunnerservice.NewDBRunnerService())
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	listenedOn := fmt.Sprintf("0.0.0.0:%s", port)
	fmt.Printf("listened at %s%s", listenedOn, path)

	err := http.ListenAndServe(listenedOn, handler)
	if err != nil {
		panic(err)
	}
}
