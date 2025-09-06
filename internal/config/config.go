package config

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/google/uuid"
)

type Config struct {
	InstanceId        string
	RedisAddr         string
	Port              string
	JwtSecret         []byte
	WindowSize        int
	WindowInterval    int
	ActiveWindowCount int
}

func Load() Config {
	instanceId := uuid.NewString()
	redisAddr := os.Getenv("REDIS_ADDR")
	if len(redisAddr) == 0 {
		redisAddr = "localhost:6379"
	}
	port := os.Getenv("PORT")
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	windowSize := parseenvInt("WINDOW_SIZE")
	windowInterval := parseenvInt("WINDOW_INTERVAL")
	activeWindowCount := parseenvInt("ACTIVE_WINDOW_COUNT")
	if activeWindowCount == 1 {
		slog.Warn("ACTIVE_WINDOW_COUNT set to 1. This may create a 'cliff edge' in slow workloads where the last ticket in the window is immediately expired if the window increments shortly after it is issued. Consider using a value above 1 to avoid this edge case")
	}

	cfg := Config{
		InstanceId:        instanceId,
		RedisAddr:         redisAddr,
		Port:              port,
		JwtSecret:         jwtSecret,
		WindowSize:        windowSize,
		WindowInterval:    windowInterval,
		ActiveWindowCount: activeWindowCount,
	}

	return cfg
}

func parseenvInt(name string) int {
	valueStr := os.Getenv(name)
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		slog.Error("Failed to parse env var", "variable", name, "error", err)
		os.Exit(1)
	}
	return value
}
