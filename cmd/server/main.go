package main

import (
	"net/http"
	"os"

	"github.com/HarrisonTCodes/session-queue/internal/redisclient"
	"github.com/HarrisonTCodes/session-queue/internal/routes"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/status", routes.HandleStatus)
	mux.HandleFunc("POST /join", routes.HandleJoin)

	redisclient.Init(os.Getenv("REDIS_ADDR"))

	http.ListenAndServe(":3000", mux)
}
