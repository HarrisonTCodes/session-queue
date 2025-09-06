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
