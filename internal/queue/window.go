package queue

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type Status string

const (
	StatusExpired Status = "expired"
	StatusWaiting Status = "waiting"
	StatusActive  Status = "active"
)

func IncrementWindow(rdb *redis.Client, ctx context.Context, size int, interval int) error {
	pipe := rdb.TxPipeline()
	pipe.IncrBy(ctx, KeyWindowEnd, int64(size))

	now := time.Now().Unix()
	newNextUpdate := now + int64(interval)
	pipe.Set(ctx, KeyNextWindowIncrement, strconv.FormatInt(newNextUpdate, 10), 0)

	_, err := pipe.Exec(ctx)

	return err
}

func GetPositionStatus(pos int64, end int64, size int, activeCount int) Status {
	if pos <= end-int64(size*activeCount) {
		return StatusExpired
	} else if pos > end {
		return StatusWaiting
	} else {
		return StatusActive
	}
}
