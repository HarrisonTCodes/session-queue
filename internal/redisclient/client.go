package redisclient

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func Init(addr string, windowSize int, windowSeconds int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	ctx := context.Background()
	err := rdb.MSetNX(ctx,
		"queue:current-position", 0,
		"queue:current-max-allowed-position", windowSize,
	).Err()
	if err != nil {
		log.Fatal(err)
	}

	go IncrWindow(rdb, ctx, windowSize, windowSeconds)

	return rdb
}

func IncrWindow(rdb *redis.Client, ctx context.Context, size int, seconds int) {
	nextIncr := time.Now().Add(time.Second * time.Duration(seconds))

	for {
		now := time.Now()
		sleepDuration := max(nextIncr.Sub(now), 0)
		time.Sleep(sleepDuration)

		vals, err := rdb.MGet(ctx, "queue:current-position", "queue:current-max-allowed-position").Result()
		if err != nil {
			log.Println("Redis error:", err)
			nextIncr = time.Now().Add(time.Second * time.Duration(3))
			continue
		}

		currentPos, _ := strconv.Atoi(vals[0].(string))
		maxPos, _ := strconv.Atoi(vals[1].(string))

		if currentPos >= maxPos {
			log.Printf("Incrementing window from %d-%d to %d-%d", maxPos-size, maxPos, maxPos, maxPos+size)
			err := rdb.IncrBy(ctx, "queue:current-max-allowed-position", int64(size)).Err()
			if err != nil {
				log.Fatal(err)
			}
			nextIncr = time.Now().Add(time.Second * time.Duration(seconds))
		} else {
			log.Printf("Skipping window increment as current position (%d) is inside current window (%d-%d)", currentPos, maxPos-size, maxPos)
			nextIncr = time.Now().Add(time.Second * time.Duration(3))
		}
	}
}
