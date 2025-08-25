package redisclient

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client

func Init(addr string) {
	Rdb = redis.NewClient(&redis.Options{
		Addr: addr,
	})

	ctx := context.Background()
	err := Rdb.MSetNX(ctx,
		"queue:current-position", 0,
		"queue:current-max-allowed-position", 0,
	).Err()
	if err != nil {
		log.Fatal(err)
	}
}
