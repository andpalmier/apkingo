package analyzer

import (
	"os"
	"testing"
)

func TestIsXAPK(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"app.xapk", true},
		{"app.apks", true},
		{"app.XAPK", true}, // Windows not case sensitive
		{"app.APKS", true},
		{"app.apk", false},
		{"app.zip", false},
		{"", false},
	}

	for _, tc := range tests {
		result := IsXAPK(tc.path)
		if result != tc.expected {
			t.Errorf("IsXAPK(%q) = %v, expected %v", tc.path, result, tc.expected)
		}
	}
}

func TestExtractAPKs(t *testing.T) {
	// Test the error case for non-existent XAPK file
	_, err := ExtractAPKs("/nonexistent/file.xapk")
	if err == nil {
		t.Error("Expected error for non-existent XAPK file")
	}
}

func TestExtractAPKs_NotAXAPK(t *testing.T) {
	// Create a regular file that's not a valid XAPK
	tmpFile, err := os.CreateTemp("", "apkingo-notxapk-*")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	if _, err := tmpFile.WriteString("not a zip file"); err != nil {
		_ = tmpFile.Close()
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}
	defer func() { _ = os.Remove(tmpFile.Name()) }()

	_, err = ExtractAPKs(tmpFile.Name())
	if err == nil {
		t.Error("Expected error for non-XAPK file")
	}
}

func TestIsAPKFile(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{"test.apk", true},
		{"TEST.APK", true},
		{"test.xapk", false},
		{"test.txt", false},
		{"test", false},
	}

	for _, tc := range tests {
		// Access the unexported function through a test wrapper
		result := isAPKFile(tc.name)
		if result != tc.expected {
			t.Errorf("isAPKFile(%q) = %v, expected %v", tc.name, result, tc.expected)
		}
	}
}
