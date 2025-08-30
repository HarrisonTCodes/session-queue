package main

import (
	"net/http"

	"github.com/HarrisonTCodes/session-queue/internal/config"
	"github.com/HarrisonTCodes/session-queue/internal/redisclient"
	"github.com/HarrisonTCodes/session-queue/internal/routes"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/status", routes.HandleStatus)
	mux.HandleFunc("POST /join", routes.HandleJoin)

	config.Load()

	redisclient.Init(config.Cfg.RedisAddr)

	http.ListenAndServe(":"+config.Cfg.Port, mux)
}
