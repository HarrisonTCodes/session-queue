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

	slog.Info("Connecting to Redis", "address", cfg.RedisAddr)
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})
	ctx := context.Background()
	queue.Init(rdb, ctx, cfg.InstanceId, cfg.RedisAddr, cfg.WindowSize, cfg.WindowInterval)

	slog.Info("Registering HTTP handlers")
	mux := http.NewServeMux()
	mux.HandleFunc("/status", routes.HandleStatus(rdb, cfg.JwtSecret, cfg.WindowSize, cfg.ActiveWindowCount))
	mux.HandleFunc("POST /join", routes.HandleJoin(rdb, cfg.JwtSecret))

	slog.Info("Running server", "port", cfg.Port)
	http.ListenAndServe(":"+cfg.Port, mux)
}
