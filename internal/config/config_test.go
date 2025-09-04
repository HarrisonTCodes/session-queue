package config

import (
	"bytes"
	"os"
	"os/exec"
	"strconv"
	"testing"
)

func TestParseenvInt_Valid(t *testing.T) {
	t.Setenv("TEST_INT", "10")
	parsed := parseenvInt("TEST_INT")
	if parsed != 10 {
		t.Fatalf("Expected 10, got %v", parsed)
	}
}

func TestParseenvInt_Invalid(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		t.Setenv("TEST_NOT_INT", "hello world")
		_ = parseenvInt("TEST_NOT_INT")
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestParseenvInt_Invalid")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()

	if err == nil {
		t.Fatalf("Expected parseenvInt to fail, no failure")
	}
}

func TestLoad(t *testing.T) {
	cfg := Config{
		InstanceId:        "id",
		RedisAddr:         "address",
		Port:              "3000",
		JwtSecret:         []byte("secret"),
		WindowSize:        10,
		WindowInterval:    10,
		ActiveWindowCount: 10,
	}

	t.Setenv("REDIS_ADDR", cfg.RedisAddr)
	t.Setenv("PORT", cfg.Port)
	t.Setenv("JWT_SECRET", string(cfg.JwtSecret))
	t.Setenv("WINDOW_SIZE", strconv.FormatInt(int64(cfg.WindowSize), 10))
	t.Setenv("WINDOW_INTERVAL", strconv.FormatInt(int64(cfg.WindowInterval), 10))
	t.Setenv("ACTIVE_WINDOW_COUNT", strconv.FormatInt(int64(cfg.ActiveWindowCount), 10))

	loadedCfg := Load()

	if cfg.RedisAddr != loadedCfg.RedisAddr ||
		cfg.Port != loadedCfg.Port ||
		cfg.WindowSize != loadedCfg.WindowSize ||
		cfg.WindowInterval != loadedCfg.WindowInterval ||
		cfg.ActiveWindowCount != loadedCfg.ActiveWindowCount ||
		!bytes.Equal(cfg.JwtSecret, loadedCfg.JwtSecret) {
		t.Fatalf("Expected %v, got %v", cfg, loadedCfg)
	}
}
