package redisclient

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func Init(addr string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	ctx := context.Background()
	err := rdb.MSetNX(ctx,
		"queue:current-position", 0,
		"queue:current-max-allowed-position", 0,
	).Err()
	if err != nil {
		log.Fatal(err)
	}

	return rdb
}
