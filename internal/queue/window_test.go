package queue

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func TestIncrementWindow(t *testing.T) {
	s, _ := miniredis.Run()
	defer s.Close()

	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{Addr: s.Addr()})
	rdb.SetNX(ctx, "queue:window-end", int64(0), 0)
	rdb.SetNX(ctx, "queue:next-window-increment", int64(0), 0)

	IncrementWindow(rdb, ctx, 10, 10)
	windowEnd, _ := rdb.Get(ctx, "queue:window-end").Result()
	if windowEnd != "10" {
		t.Fatalf("Expected window-end to be 10, got %v", windowEnd)
	}

	nextWindowIncrement, _ := rdb.Get(ctx, "queue:window-end").Result()
	if windowEnd == "0" {
		t.Fatalf("Expected next-window-increment to have changed, got %v", nextWindowIncrement)
	}
}
