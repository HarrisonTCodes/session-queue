package redisclient

import (
	"context"
	"log"
	"strconv"
	"time"

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
		log.Fatal(err)
	}

	go incrWindow(rdb, ctx, instanceId, windowSize, windowInterval)

	return rdb
}

func incrWindow(rdb *redis.Client, ctx context.Context, instanceId string, size int, interval int) {
	checkDuration := time.Second * 3
	leaderDuration := time.Second * 5

	for {
		isLeader, err := ensureElectedLeader(rdb, ctx, instanceId, leaderDuration)
		if err != nil {
			time.Sleep(checkDuration)
			continue
		}

		if !isLeader {
			log.Println("Self is not leader")
			time.Sleep(checkDuration)
			continue
		}

		vals, err := rdb.MGet(ctx, "queue:current-position", "queue:window-end", "queue:next-window-increment").Result()
		if err != nil {
			log.Println("Redis error during retrieval of queue values:", err)
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
			log.Printf("Incrementing window from %d-%d to %d-%d", windowEnd-size, windowEnd, windowEnd, windowEnd+size)

			pipe := rdb.TxPipeline()
			pipe.IncrBy(ctx, "queue:window-end", int64(size))
			newNextUpdate := now + int64(interval)
			pipe.Set(ctx, "queue:next-window-increment", strconv.FormatInt(newNextUpdate, 10), 0)
			_, err := pipe.Exec(ctx)
			if err != nil {
				log.Println("Redis error during incrementing window", err)
				time.Sleep(checkDuration)
				continue
			}
		} else {
			log.Printf("Skipping window increment as current position (%d) is inside window (%d-%d)", currentPos, windowEnd-size, windowEnd)
		}

		time.Sleep(checkDuration)
	}
}
