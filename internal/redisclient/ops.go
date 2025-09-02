package redisclient

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func ensureElectedLeader(rdb *redis.Client, ctx context.Context, instanceId string, leaderDuration time.Duration) (bool, error) {
	leaderId, err := rdb.Get(ctx, "queue:leader-id").Result()
	if err != nil && err != redis.Nil {
		log.Println("Redis error during leader retrieval:", err)
		return false, err
	}

	var isLeader bool
	if err == redis.Nil {
		log.Println("Electing self as leader")
		setLeader, err := rdb.SetNX(ctx, "queue:leader-id", instanceId, leaderDuration).Result()
		if err != nil {
			log.Println("Redis error during leader election:", err)
			return false, err
		}
		isLeader = setLeader
	} else {
		isLeader = leaderId == instanceId
	}

	if isLeader {
		log.Println("Self is leader, extending expiry")
		err = rdb.Expire(ctx, "queue:leader-id", leaderDuration).Err()
		if err != nil {
			log.Println("Redis error during leadership extension:", err)
			return false, err
		}
	}

	return isLeader, nil
}
