package config

import (
	"os"
	"testing"
)

func TestGetAPIKey(t *testing.T) {
	// Test case 1: Flag value provided
	val := getAPIKey("flag-key", "TEST_ENV_VAR", "msg")
	if val != "flag-key" {
		t.Errorf("Expected 'flag-key', got '%s'", val)
	}

	// Test case 2: Env var provided
	os.Setenv("TEST_ENV_VAR", "env-key")
	defer os.Unsetenv("TEST_ENV_VAR")
	val = getAPIKey("", "TEST_ENV_VAR", "msg")
	if val != "env-key" {
		t.Errorf("Expected 'env-key', got '%s'", val)
	}

	// Test case 3: Neither provided
	os.Unsetenv("TEST_ENV_VAR")
	val = getAPIKey("", "TEST_ENV_VAR", "msg")
	if val != "" {
		t.Errorf("Expected empty string, got '%s'", val)
	}
}
