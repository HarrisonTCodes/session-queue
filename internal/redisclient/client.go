package redisclient

import "github.com/redis/go-redis/v9"

var Rdb *redis.Client

func Init(addr string) {
	Rdb = redis.NewClient(&redis.Options{
		Addr: addr,
	})
}
