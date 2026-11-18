package vt

import (
	"encoding/json"
	"testing"
)

func TestParseVTAnalysStats(t *testing.T) {
	got := parseVTAnalysStats(map[string]interface{}{
		"harmless":          json.Number("10"),
		"type-unsupported":  json.Number("2"),
		"suspicious":        json.Number("3"),
		"confirmed-timeout": json.Number("4"),
		"timeout":           json.Number("5"),
		"failure":           json.Number("6"),
		"malicious":         json.Number("7"),
		"undetected":        json.Number("8"),
	})

	if got.Harmless != 10 || got.TypeUnsupported != 2 || got.Suspicious != 3 ||
		got.ConfirmedTimeout != 4 || got.Timeout != 5 || got.Failure != 6 ||
		got.Malicious != 7 || got.Undetected != 8 {
		t.Fatalf("Unexpected stats: %+v", got)
	}
}

func TestParseAndroguardInfo(t *testing.T) {
	got := parseAndroguardInfo(map[string]interface{}{
		"Package":     "com.example",
		"Activities":  []interface{}{"MainActivity", 42, "SettingsActivity"},
		"Services":    []interface{}{"SyncService"},
		"certificate": map[string]interface{}{"issuer": "example"},
		"permission_details": map[string]interface{}{
			"android.permission.CAMERA":   map[string]interface{}{"permission_type": "dangerous"},
			"android.permission.INTERNET": map[string]interface{}{"permission_type": "normal"},
		},
	})

	if got.Package != "com.example" {
		t.Errorf("Unexpected package: %q", got.Package)
	}
	if len(got.Activities) != 2 || got.Activities[0] != "MainActivity" || got.Activities[1] != "SettingsActivity" {
		t.Errorf("Unexpected activities: %v", got.Activities)
	}
	if len(got.DangerPerm) != 1 || got.DangerPerm[0] != "android.permission.CAMERA" {
		t.Errorf("Unexpected dangerous permissions: %v", got.DangerPerm)
	}
}
