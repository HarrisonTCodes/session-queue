package main

import (
	"log"
	"net/http"

	"github.com/HarrisonTCodes/session-queue/internal/config"
	"github.com/HarrisonTCodes/session-queue/internal/redisclient"
	"github.com/HarrisonTCodes/session-queue/internal/routes"
)

func main() {
	log.Println("Loading config")
	cfg := config.Load()

	log.Println("Connecting to Redis")
	rdb := redisclient.Init(cfg.RedisAddr, cfg.WindowSize, cfg.WindowInterval)

	log.Println("Registering handlers")
	mux := http.NewServeMux()
	mux.HandleFunc("/status", routes.HandleStatus(rdb))
	mux.HandleFunc("POST /join", routes.HandleJoin(rdb))

	log.Printf("Server running on port %s", cfg.Port)
	http.ListenAndServe(":"+cfg.Port, mux)
}
