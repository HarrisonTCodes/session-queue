package main

import (
	"log"
	"net/http"

	"github.com/HarrisonTCodes/session-queue/internal/config"
	"github.com/HarrisonTCodes/session-queue/internal/redisclient"
	"github.com/HarrisonTCodes/session-queue/internal/routes"
)

func main() {
	log.Println("Registering handlers")
	mux := http.NewServeMux()
	mux.HandleFunc("/status", routes.HandleStatus)
	mux.HandleFunc("POST /join", routes.HandleJoin)

	log.Println("Loading config")
	config.Load()

	log.Println("Connecting to Redis")
	redisclient.Init(config.Cfg.RedisAddr)

	log.Printf("Server running on port %s", config.Cfg.Port)
	http.ListenAndServe(":"+config.Cfg.Port, mux)
}
