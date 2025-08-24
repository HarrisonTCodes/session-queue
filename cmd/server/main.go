package main

import (
	"context"
	"log"
	"net/http"

	"github.com/HarrisonTCodes/session-queue/internal/redisclient"
	"github.com/HarrisonTCodes/session-queue/internal/routes"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/status", routes.HandleStatus)
	mux.HandleFunc("POST /join", routes.HandleJoin)

	redisclient.Init("localhost:6379")
	ctx := context.Background()
	err := redisclient.Rdb.Set(ctx, "foo", 0, 0).Err()
	if err != nil {
		log.Fatal(err)
	}

	http.ListenAndServe(":3000", mux)
}
