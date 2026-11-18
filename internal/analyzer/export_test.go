package analyzer_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/andpalmier/apkingo/v3/internal/analyzer"
)

func TestExportJSONAtomicAndPrivate(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "report.json")
	if err := os.WriteFile(path, []byte("old"), 0644); err != nil {
		t.Fatalf("Failed to seed report: %v", err)
	}

	app := analyzer.AndroidApp{Name: "example", PackageName: "com.example"}
	if err := app.ExportJSON(path); err != nil {
		t.Fatalf("ExportJSON failed: %v", err)
	}

	stat, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Failed to stat report: %v", err)
	}
	if got := stat.Mode().Perm(); got != 0600 {
		t.Errorf("Expected report mode 0600, got %04o", got)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read report: %v", err)
	}
	var decoded analyzer.AndroidApp
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Report is not valid JSON: %v", err)
	}
	if decoded.Name != app.Name || decoded.PackageName != app.PackageName {
		t.Errorf("Unexpected report contents: %+v", decoded)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("Failed to read report directory: %v", err)
	}
	if len(entries) != 1 {
		t.Errorf("Expected no leftover temporary files, got %v", entries)
	}
}
