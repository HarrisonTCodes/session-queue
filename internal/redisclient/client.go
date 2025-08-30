package redisclient

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func Init(addr string, windowSize int, windowInterval int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	ctx := context.Background()
	err := rdb.MSetNX(ctx,
		"queue:current-position", 0,
		"queue:window-end", windowSize,
	).Err()
	if err != nil {
		log.Fatal(err)
	}

	go IncrWindow(rdb, ctx, windowSize, windowInterval)

	return rdb
}

func IncrWindow(rdb *redis.Client, ctx context.Context, size int, interval int) {
	checkDuration := time.Second * time.Duration(3)
	intervalDuration := time.Second * time.Duration(interval)
	nextIncr := time.Now().Add(intervalDuration)

	for {
		now := time.Now()
		sleepDuration := max(nextIncr.Sub(now), 0)
		time.Sleep(sleepDuration)

		vals, err := rdb.MGet(ctx, "queue:current-position", "queue:window-end").Result()
		if err != nil {
			log.Println("Redis error:", err)
			nextIncr = time.Now().Add(checkDuration)
			continue
		}

		currentPos, _ := strconv.Atoi(vals[0].(string))
		end, _ := strconv.Atoi(vals[1].(string))

		if currentPos >= end {
			log.Printf("Incrementing window from %d-%d to %d-%d", end-size, end, end, end+size)
			err := rdb.IncrBy(ctx, "queue:window-end", int64(size)).Err()
			if err != nil {
				log.Fatal(err)
			}
			nextIncr = time.Now().Add(intervalDuration)
		} else {
			log.Printf("Skipping window increment as current position (%d) is inside current window (%d-%d)", currentPos, end-size, end)
			nextIncr = time.Now().Add(checkDuration)
		}
	}
}
