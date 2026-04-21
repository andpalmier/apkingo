package analyzer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindAPKs(t *testing.T) {
	// Create a temporary directory with test files
	tmpDir, err := os.MkdirTemp("", "apkingo-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create test files
	testFiles := []string{
		"test1.apk",
		"test2.xapk",
		"test3.apks",
		"test4.txt", // should be ignored
		"test5.APK", // should be found (case insensitive check - currently not)
	}

	for _, name := range testFiles {
		if _, err := os.Create(filepath.Join(tmpDir, name)); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	// Create a subdirectory (should be ignored)
	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}

	// Test findAPKs
	paths, err := findAPKs(tmpDir)
	if err != nil {
		t.Fatalf("findAPKs failed: %v", err)
	}

	// Should find 3 APK files (.apk, .xapk, .apks)
	expected := 3
	if len(paths) != expected {
		t.Errorf("Expected %d APK files, got %d: %v", expected, len(paths), paths)
	}
}

func TestAnalyzeDirectory_EmptyDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "apkingo-empty-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	_, _, err = AnalyzeDirectory(tmpDir, "us", "", "", false)
	if err == nil {
		t.Error("Expected error for empty directory")
	}
}

func TestAnalyzeDirectory_NotADirectory(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "apkingo-notdir-*")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}
	defer func() { _ = os.Remove(tmpFile.Name()) }()

	_, _, err = AnalyzeDirectory(tmpFile.Name(), "us", "", "", false)
	if err == nil {
		t.Error("Expected error for non-directory path")
	}
}

func TestAnalyzeDirectory_NonExistent(t *testing.T) {
	_, _, err := AnalyzeDirectory("/nonexistent/path/12345", "us", "", "", false)
	if err == nil {
		t.Error("Expected error for non-existent path")
	}
}
