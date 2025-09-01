package config

import (
	"testing"
)

func TestParseenvInt(t *testing.T) {
	t.Setenv("TEST_INT", "10")
	parsed := parseenvInt("TEST_INT")
	if parsed != 10 {
		t.Errorf("Expected 10, got %v", parsed)
	}
}
