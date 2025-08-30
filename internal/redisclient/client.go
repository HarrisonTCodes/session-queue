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

	windowStartStr, err := rdb.Get(ctx, "queue:current-max-allowed-position").Result()
	if err != nil {
		log.Fatalln(err)
	}
	windowStart, _ := strconv.Atoi(windowStartStr)

	go IncrWindow(rdb, ctx, windowStart, windowSize, windowSeconds)

	return rdb
}

func IncrWindow(rdb *redis.Client, ctx context.Context, start int, size int, seconds int) {
	ticker := time.NewTicker(time.Second * time.Duration(seconds))

	defer ticker.Stop()

	for range ticker.C {
		log.Printf("Incrementing window from %d to %d", start, start+size)
		start += size
		err := rdb.IncrBy(ctx, "queue:current-max-allowed-position", int64(size)).Err()
		if err != nil {
			log.Fatal(err)
		}
	}
}
