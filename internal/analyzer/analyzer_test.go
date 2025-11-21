package analyzer_test

import (
	"testing"

	"github.com/andpalmier/apkingo/internal/analyzer"
)

func TestProcessAPK(t *testing.T) {
	// Use F-Droid.apk for testing
	apkPath := "../../test/F-Droid.apk"

	app := analyzer.AndroidApp{}

	// Test with no API keys
	err := app.ProcessAPK(apkPath, "us", "", "")
	if err != nil {
		t.Fatalf("ProcessAPK failed: %v", err)
	}

	// Verify some basic info
	if app.PackageName != "org.fdroid.fdroid" {
		t.Errorf("Expected package name org.fdroid.fdroid, got %s", app.PackageName)
	}

	if app.Hashes.Md5 != "df1373f9fd535abddd86d3a5a9c87bbe" {
		t.Errorf("Expected MD5 df1373f9fd535abddd86d3a5a9c87bbe, got %s", app.Hashes.Md5)
	}

	if len(app.Permissions) == 0 {
		t.Error("Expected permissions to be found, got 0")
	}

	if app.Certificate.Serial == "" {
		t.Error("Expected certificate serial to be found")
	}
}

func TestProcessAPK_FileNotFound(t *testing.T) {
	app := analyzer.AndroidApp{}
	err := app.ProcessAPK("nonexistent.apk", "us", "", "")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}
