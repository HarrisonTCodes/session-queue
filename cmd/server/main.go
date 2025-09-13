package main

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/HarrisonTCodes/session-queue/internal/config"
	"github.com/HarrisonTCodes/session-queue/internal/queue"
	"github.com/HarrisonTCodes/session-queue/internal/routes"
	"github.com/redis/go-redis/v9"
)

func main() {
	slog.Info("Loading config")
	cfg := config.Load()

	slog.Info("Creating Redis client", "address", cfg.RedisAddr)
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	slog.Info("Registering HTTP handlers")
	mux := http.NewServeMux()
	mux.HandleFunc("/livez", routes.HandleLivez)
	mux.HandleFunc("/status", routes.HandleStatus(rdb, cfg.JwtSecret, cfg.WindowSize, cfg.ActiveWindowCount))
	mux.HandleFunc("POST /join", routes.HandleJoin(rdb, cfg.JwtSecret))

	slog.Info("Running server", "port", cfg.Port)
	go http.ListenAndServe(":"+cfg.Port, mux)

	slog.Info("Setting up queue in Redis")
	ctx := context.Background()
	queue.Init(rdb, ctx, cfg.InstanceId, cfg.RedisAddr, cfg.WindowSize, cfg.WindowInterval)

	select {}
}
