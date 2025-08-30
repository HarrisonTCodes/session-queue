package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	RedisAddr     string
	Port          int
	WindowSize    int
	WindowSeconds int
}

var Cfg Config

func Load() {
	redisAddr := os.Getenv("REDIS_ADDR")
	port := parseenvInt("PORT")
	windowSize := parseenvInt("WINDOW_SIZE")
	windowSeconds := parseenvInt("WINDOW_SECONDS")

	Cfg = Config{
		RedisAddr:     redisAddr,
		Port:          port,
		WindowSize:    windowSize,
		WindowSeconds: windowSeconds,
	}
}

func parseenvInt(name string) int {
	valueStr := os.Getenv(name)
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Fatalf("Failed to parse %s env var: %v", name, err)
	}
	return value
}
