package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	RedisAddr     string
	Port          string
	WindowSize    int
	WindowSeconds int
}

func Load() Config {
	redisAddr := os.Getenv("REDIS_ADDR")
	port := os.Getenv("PORT")
	windowSize := parseenvInt("WINDOW_SIZE")
	windowSeconds := parseenvInt("WINDOW_SECONDS")

	cfg := Config{
		RedisAddr:     redisAddr,
		Port:          port,
		WindowSize:    windowSize,
		WindowSeconds: windowSeconds,
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
