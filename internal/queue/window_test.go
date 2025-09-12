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
	rdb.SetNX(ctx, KeyWindowEnd, int64(0), 0)
	rdb.SetNX(ctx, KeyNextWindowIncrement, int64(0), 0)

	IncrementWindow(rdb, ctx, 10, 10)
	windowEnd, _ := rdb.Get(ctx, KeyWindowEnd).Result()
	if windowEnd != "10" {
		t.Fatalf("Expected window-end to be 10, got %v", windowEnd)
	}

	nextWindowIncrement, _ := rdb.Get(ctx, KeyWindowEnd).Result()
	if windowEnd == "0" {
		t.Fatalf("Expected next-window-increment to have changed, got %v", nextWindowIncrement)
	}
}

func TestGetPositionStatus(t *testing.T) {
	if status := GetPositionStatus(1, 20, 10, 1); status != StatusExpired {
		t.Fatalf("Expected status to be Expired, got %v", status)
	}

	if status := GetPositionStatus(15, 20, 10, 1); status != StatusActive {
		t.Fatalf("Expected status to be Active, got %v", status)
	}

	if status := GetPositionStatus(30, 20, 10, 1); status != StatusWaiting {
		t.Fatalf("Expected status to be Waiting, got %v", status)
	}
}
