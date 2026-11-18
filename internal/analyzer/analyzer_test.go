package analyzer_test

import (
	"slices"
	"testing"

	"github.com/andpalmier/apkingo/v3/internal/analyzer"
)

func TestProcessAPK(t *testing.T) {
	// Use F-Droid.apk for testing
	apkPath := "../../test/F-Droid.apk"

	app := analyzer.AndroidApp{}

	// Test with no API keys
	err := app.ProcessAPK(apkPath, "us", "", "", "", true)
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
	if !app.Certificate.Verified {
		t.Errorf("Expected certificate to verify, got error: %s", app.Certificate.VerificationError)
	}
	if app.Certificate.SigningScheme != "v3" {
		t.Errorf("Expected v3 signing scheme, got %s", app.Certificate.SigningScheme)
	}
	if app.Certificate.SignerCount != 1 {
		t.Errorf("Expected one signer, got %d", app.Certificate.SignerCount)
	}
	if app.Certificate.Sha256 == "" {
		t.Error("Expected certificate SHA256 to be found")
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

func TestSetSecurityFlags(t *testing.T) {
	tests := []struct {
		name                 string
		apkPath              string
		wantAllowBackup      bool
		wantCleartext        *bool
		wantNetworkConfigSet bool
	}{
		{
			name:                 "F-Droid",
			apkPath:              "../../test/F-Droid.apk",
			wantAllowBackup:      true,
			wantNetworkConfigSet: true,
		},
		{
			name:            "Snapseed",
			apkPath:         "../../test/snapseed.apk",
			wantAllowBackup: true,
			wantCleartext:   boolPtr(false),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := analyzer.AndroidApp{}
			if err := app.SetSecurityFlags(tt.apkPath); err != nil {
				t.Fatalf("SetSecurityFlags failed: %v", err)
			}
			if app.AllowBackup == nil || *app.AllowBackup != tt.wantAllowBackup {
				t.Errorf("AllowBackup = %v, want %t", app.AllowBackup, tt.wantAllowBackup)
			}
			if tt.wantCleartext == nil {
				if app.UsesCleartextTraffic != nil {
					t.Errorf("UsesCleartextTraffic = %v, want not set", *app.UsesCleartextTraffic)
				}
			} else if app.UsesCleartextTraffic == nil || *app.UsesCleartextTraffic != *tt.wantCleartext {
				t.Errorf("UsesCleartextTraffic = %v, want %t", app.UsesCleartextTraffic, *tt.wantCleartext)
			}
			if got := app.NetworkSecurityConfig != ""; got != tt.wantNetworkConfigSet {
				t.Errorf("NetworkSecurityConfig set = %t, want %t", got, tt.wantNetworkConfigSet)
			}
		})
	}
}

func TestSetSecurityFlagsComponents(t *testing.T) {
	tests := []struct {
		name           string
		apkPath        string
		wantActivities []string
		wantServices   []string
		wantProviders  []string
		wantReceivers  []string
	}{
		{
			name:           "F-Droid",
			apkPath:        "../../test/F-Droid.apk",
			wantActivities: []string{"org.fdroid.fdroid.views.main.MainActivity"},
			wantServices:   []string{"org.fdroid.fdroid.UpdateService"},
			wantProviders:  []string{"org.fdroid.fdroid.data.AppProvider"},
			wantReceivers:  []string{"org.fdroid.fdroid.receiver.StartupReceiver"},
		},
		{
			name:    "Snapseed",
			apkPath: "../../test/snapseed.apk",
			wantActivities: []string{
				"com.google.android.apps.snapseed.activities.edit.MainActivity",
				"com.google.android.apps.snapseed.MainActivity",
			},
			wantServices:  []string{"com.google.android.apps.snapseed.core.SaveFilteredImageLocalService"},
			wantProviders: []string{"android.support.v4.content.FileProvider"},
			wantReceivers: []string{"com.google.android.libraries.social.silentfeedback.impl.SilentFeedbackReceiver"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := analyzer.AndroidApp{}
			if err := app.SetSecurityFlags(tt.apkPath); err != nil {
				t.Fatalf("SetSecurityFlags failed: %v", err)
			}

			checks := []struct {
				label string
				got   []string
				want  []string
			}{
				{"Activities", app.Activities, tt.wantActivities},
				{"Services", app.Services, tt.wantServices},
				{"Providers", app.Providers, tt.wantProviders},
				{"Receivers", app.Receivers, tt.wantReceivers},
			}
			for _, check := range checks {
				for _, want := range check.want {
					if !slices.Contains(check.got, want) {
						t.Errorf("%s missing %q: %v", check.label, want, check.got)
					}
				}
			}
		})
	}
}

func boolPtr(v bool) *bool {
	return &v
}

func TestProcessAPK_FileNotFound(t *testing.T) {
	app := analyzer.AndroidApp{}
	err := app.ProcessAPK("nonexistent.apk", "us", "", "", "", false)
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}
