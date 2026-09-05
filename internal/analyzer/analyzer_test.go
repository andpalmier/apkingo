package analyzer_test

import (
	"testing"

	"github.com/andpalmier/apkingo/v2/internal/analyzer"
)

func TestProcessAPK(t *testing.T) {
	// Use F-Droid.apk for testing
	apkPath := "../../test/F-Droid.apk"

	app := analyzer.AndroidApp{}

	// Test with no API keys
	err := app.ProcessAPK(apkPath, "us", "", "", "", false)
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

	// Verify version fields
	if app.VersionName == "" {
		t.Error("Expected VersionName to be non-empty")
	}
	if app.VersionCode == 0 {
		t.Error("Expected VersionCode to be non-zero")
	}
	if app.VersionCode < 0 {
		t.Errorf("Expected positive VersionCode, got %d", app.VersionCode)
	}

	// F-Droid is pure Java/Kotlin — no native libs expected
	if app.Architectures == nil {
		t.Error("Expected Architectures slice to be initialized (can be empty)")
	}
}

func TestProcessAPK_FileNotFound(t *testing.T) {
	app := analyzer.AndroidApp{}
	err := app.ProcessAPK("nonexistent.apk", "us", "", "", "", false)
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}
