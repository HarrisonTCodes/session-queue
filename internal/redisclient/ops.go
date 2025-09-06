package redisclient

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func ensureElectedLeader(rdb *redis.Client, ctx context.Context, instanceId string, leaderDuration time.Duration) (bool, error) {
	leaderId, err := rdb.Get(ctx, "queue:leader-id").Result()
	if err != nil && err != redis.Nil {
		slog.Error("Redis error during leader retrieval", "error", err)
		return false, err
	}

	var isLeader bool
	if err == redis.Nil {
		slog.Info("Electing self as leader")
		setLeader, err := rdb.SetNX(ctx, "queue:leader-id", instanceId, leaderDuration).Result()
		if err != nil {
			slog.Error("Redis error during leader election", "error", err)
			return false, err
		}
		isLeader = setLeader
	} else {
		isLeader = leaderId == instanceId
	}

	if isLeader {
		slog.Info("Self is leader, extending expiry")
		err = rdb.Expire(ctx, "queue:leader-id", leaderDuration).Err()
		if err != nil {
			slog.Error("Redis error during leadership extension", "error", err)
			return false, err
		}
	}

	return isLeader, nil
}

func incrementWindow(rdb *redis.Client, ctx context.Context, size int, interval int) error {
	pipe := rdb.TxPipeline()
	pipe.IncrBy(ctx, "queue:window-end", int64(size))

	now := time.Now().Unix()
	newNextUpdate := now + int64(interval)
	pipe.Set(ctx, "queue:next-window-increment", strconv.FormatInt(newNextUpdate, 10), 0)

	_, err := pipe.Exec(ctx)

	return err
}
