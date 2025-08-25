package main

import (
	"net/http"

	"github.com/HarrisonTCodes/session-queue/internal/redisclient"
	"github.com/HarrisonTCodes/session-queue/internal/routes"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/status", routes.HandleStatus)
	mux.HandleFunc("POST /join", routes.HandleJoin)

	redisclient.Init("localhost:6379")

	http.ListenAndServe(":3000", mux)
}
