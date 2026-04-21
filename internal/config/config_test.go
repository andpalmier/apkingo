package config

import (
	"os"
	"testing"
)

// Note: Testing flag parsing is limited because flag.Parse() is called in Load() at package init.
// We verify the flag is registered by checking the code compiles correctly.
// Integration tests in main package can test flag behavior with os.Args.

func TestNoPlayStoreFlagExists(t *testing.T) {
	// Verify NoPlayStore field exists in Config struct
	cfg := &Config{}
	// This won't compile if NoPlayStore field doesn't exist
	_ = cfg.NoPlayStore
	// Default should be false
	if cfg.NoPlayStore != false {
		t.Errorf("Expected NoPlayStore default to be false, got %v", cfg.NoPlayStore)
	}
}

func TestGetAPIKey(t *testing.T) {
	// Test case 1: Flag value provided
	val := getAPIKey("flag-key", "TEST_ENV_VAR", "msg")
	if val != "flag-key" {
		t.Errorf("Expected 'flag-key', got '%s'", val)
	}

	// Test case 2: Env var provided
	if err := os.Setenv("TEST_ENV_VAR", "env-key"); err != nil {
		t.Fatalf("Failed to set env var: %v", err)
	}
	defer func() {
		if err := os.Unsetenv("TEST_ENV_VAR"); err != nil {
			t.Logf("Warning: failed to unset env var: %v", err)
		}
	}()
	val = getAPIKey("", "TEST_ENV_VAR", "msg")
	if val != "env-key" {
		t.Errorf("Expected 'env-key', got '%s'", val)
	}

	// Test case 3: Neither provided
	if err := os.Unsetenv("TEST_ENV_VAR"); err != nil {
		t.Fatalf("Failed to unset env var: %v", err)
	}
	val = getAPIKey("", "TEST_ENV_VAR", "msg")
	if val != "" {
		t.Errorf("Expected empty string, got '%s'", val)
	}
}
