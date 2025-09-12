package redisclient

import (
	"context"
	"log/slog"
	"os"

	"github.com/HarrisonTCodes/session-queue/internal/queue"
	"github.com/redis/go-redis/v9"
)

func Init(instanceId string, addr string, windowSize int, windowInterval int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	ctx := context.Background()
	err := rdb.MSetNX(ctx,
		"queue:current-position", 0,
		"queue:window-end", windowSize,
		"queue:next-window-increment", windowInterval,
	).Err()
	if err != nil {
		slog.Error("Redis error during initialisation of key-value pairs", "error", err)
		os.Exit(1)
	}

	go queue.WatchQueue(rdb, ctx, instanceId, windowSize, windowInterval)

	return rdb
}
