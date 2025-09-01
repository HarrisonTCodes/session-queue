package config

import (
	"os"
	"os/exec"
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
