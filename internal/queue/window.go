package queue

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func IncrementWindow(rdb *redis.Client, ctx context.Context, size int, interval int) error {
	pipe := rdb.TxPipeline()
	pipe.IncrBy(ctx, "queue:window-end", int64(size))

	now := time.Now().Unix()
	newNextUpdate := now + int64(interval)
	pipe.Set(ctx, "queue:next-window-increment", strconv.FormatInt(newNextUpdate, 10), 0)

	_, err := pipe.Exec(ctx)

	return err
}
