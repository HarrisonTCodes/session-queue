package config

import (
	"log"
	"os"
	"strconv"

	"github.com/google/uuid"
)

type Config struct {
	InstanceId        string
	RedisAddr         string
	Port              string
	WindowSize        int
	WindowInterval    int
	ActiveWindowCount int
}

func Load() Config {
	instanceId := uuid.NewString()
	redisAddr := os.Getenv("REDIS_ADDR")
	port := os.Getenv("PORT")
	windowSize := parseenvInt("WINDOW_SIZE")
	windowInterval := parseenvInt("WINDOW_INTERVAL")
	activeWindowCount := parseenvInt("ACTIVE_WINDOW_COUNT")

	cfg := Config{
		InstanceId:        instanceId,
		RedisAddr:         redisAddr,
		Port:              port,
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
		log.Fatalf("Failed to parse %s env var: %v", name, err)
	}
	return value
}
