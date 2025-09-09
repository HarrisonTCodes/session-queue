package queue

import (
	"context"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

func EnsureElectedLeader(rdb *redis.Client, ctx context.Context, instanceId string, leaderDuration time.Duration) (bool, error) {
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
