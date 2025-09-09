package redisclient

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

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

	go watchQueue(rdb, ctx, instanceId, windowSize, windowInterval)

	return rdb
}

func watchQueue(rdb *redis.Client, ctx context.Context, instanceId string, windowSize int, interval int) {
	checkDuration := time.Second * 3
	leaderDuration := time.Second * 5

	for {
		isLeader, err := queue.EnsureElectedLeader(rdb, ctx, instanceId, leaderDuration)
		if err != nil {
			time.Sleep(checkDuration)
			continue
		}

		if !isLeader {
			slog.Info("Self is not leader, skipping")
			time.Sleep(checkDuration)
			continue
		}

		vals, err := rdb.MGet(ctx, "queue:current-position", "queue:window-end", "queue:next-window-increment").Result()
		if err != nil {
			slog.Error("Redis error during retrieval of queue values", "error", err)
			time.Sleep(checkDuration)
			continue
		}

		currentPos, _ := strconv.Atoi(vals[0].(string))
		windowEnd, _ := strconv.Atoi(vals[1].(string))
		nextUpdate, _ := strconv.ParseInt(vals[2].(string), 10, 64)

		now := time.Now().Unix()
		sleepDuration := time.Second * time.Duration(max(nextUpdate-now, 0))
		if sleepDuration > 0 {
			time.Sleep(sleepDuration)
			continue
		}

		if currentPos > windowEnd {
			slog.Info("Incrementing window", "previous", fmt.Sprintf("%d-%d", windowEnd-windowSize, windowEnd), "new", fmt.Sprintf("%d-%d", windowEnd, windowEnd+windowSize))

			err = queue.IncrementWindow(rdb, ctx, windowSize, interval)
			if err != nil {
				slog.Error("Redis error during incrementing window", "error", err)
				time.Sleep(checkDuration)
				continue
			}
		} else {
			slog.Info("Skipping window increment as current position is inside window", "position", currentPos, "window", fmt.Sprintf("%d-%d", windowEnd-windowSize, windowEnd))
		}

		time.Sleep(checkDuration)
	}
}
