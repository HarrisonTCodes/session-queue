package redisclient

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func TestElectSelfAsLeader(t *testing.T) {
	s, _ := miniredis.Run()
	defer s.Close()

	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{Addr: s.Addr()})
	ensureElectedLeader(rdb, ctx, "id", time.Second*10)

	leaderId, _ := rdb.Get(ctx, "queue:leader-id").Result()
	if leaderId != "id" {
		t.Fatalf("Expected leader to be id, got %v", leaderId)
	}
}

func TestLeaderAlreadyElected(t *testing.T) {
	s, _ := miniredis.Run()
	defer s.Close()

	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{Addr: s.Addr()})
	rdb.SetNX(ctx, "queue:leader-id", "other-id", 0)
	ensureElectedLeader(rdb, ctx, "id", time.Second*10)

	leaderId, _ := rdb.Get(ctx, "queue:leader-id").Result()
	if leaderId != "other-id" {
		t.Fatalf("Expected leader to be other-id, got %v", leaderId)
	}
}
