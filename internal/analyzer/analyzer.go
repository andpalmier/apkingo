package analyzer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/andpalmier/apkingo/v3/internal/koodous"
	"github.com/andpalmier/apkingo/v3/internal/vt"
	"github.com/shogo82148/androidbinary/apk"
)

// AndroidApp represents information extracted from an APK file
type AndroidApp struct {
	Name                  string               `json:"name"`
	NameLocale            string               `json:"name-locale,omitempty"` // localized name from -locale flag
	Locale                string               `json:"locale,omitempty"`      // locale requested (e.g., "zh-CN")
	PackageName           string               `json:"package-name"`
	VersionName           string               `json:"version-name"`
	VersionCode           int32                `json:"version-code"`
	Architectures         []string             `json:"architectures,omitempty"`
	MainActivity          string               `json:"main-activity"`
	MinimumSDK            int32                `json:"minimum-sdk"`
	TargetSDK             int32                `json:"target-sdk"`
	Hashes                Hashes               `json:"hashes"`
	Permissions           []string             `json:"permissions"`
	PermissionClasses     map[string]string    `json:"permission-classes,omitempty"`
	Metadata              []Metadata           `json:"metadata"`
	Debuggable            *bool                `json:"debuggable,omitempty"`
	AllowBackup           *bool                `json:"allow-backup,omitempty"`
	TestOnly              *bool                `json:"test-only,omitempty"`
	UsesCleartextTraffic  *bool                `json:"uses-cleartext-traffic,omitempty"`
	NetworkSecurityConfig string               `json:"network-security-config,omitempty"`
	Certificate           CertificateInfo      `json:"certificate"`
	PlayStore             *PlayStoreInfo       `json:"playstore,omitempty"`
	Koodous               *koodous.KoodousInfo `json:"koodous,omitempty"`
	VirusTotal            *vt.VirusTotalInfo   `json:"virustotal,omitempty"`
	NoPlayStore           bool                 `json:"-"` // internal use only - skip Play Store
	Errors                AnalysisErrors       `json:"-"` // internal use only
}

// AnalysisErrors holds errors encountered during analysis
type AnalysisErrors struct {
	General   error
	Security  error
	Cert      error
	PlayStore error
	Koodous   error
	VT        error
}

// Hashes represents hash values
type Hashes struct {
	Md5    string `json:"md5"`
	Sha1   string `json:"sha1"`
	Sha256 string `json:"sha256"`
}

// Metadata represents metadata
type Metadata struct {
	Name  string `json:"name"`
	Value string `json:"value,omitempty"`
}

// ProcessAPK orchestrates the APK analysis.
// country is the Play Store country code for localized store results.
// locale is the language tag for localized app name extraction (e.g., "en", "zh-CN").
// noPlayStore when true, skips Play Store API calls for offline analysis.
func (app *AndroidApp) ProcessAPK(apkPath, country, locale, vtAPIKey, koodousAPI string, noPlayStore bool) error {
	pkg, err := apk.OpenFile(apkPath)
	if err != nil {
		return fmt.Errorf("error loading APK: %w", err)
	}
	defer func() {
		if err := pkg.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to close APK: %v\n", err)
		}
	}()

	if err = app.SetGeneralInfo(pkg, locale); err != nil {
		app.Errors.General = err
	}

	if err = app.SetSecurityFlags(apkPath); err != nil {
		app.Errors.Security = err
	}

	app.SetArchitectures(apkPath)

	if err = app.SetHashes(apkPath); err != nil {
		return fmt.Errorf("error setting hashes: %w", err)
	}

	if err = app.SetCertInfo(apkPath); err != nil {
		app.Errors.Cert = err
	}

	// Track if user requested no Play Store
	app.NoPlayStore = noPlayStore

	// Run external API calls concurrently for better performance
	var wg sync.WaitGroup

	// PlayStore - skip if noPlayStore flag is set
	if !noPlayStore {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := app.SetPlayStoreInfo(country); err != nil {
				app.Errors.PlayStore = err
				// Error is stored for reporting and will be displayed in output
			}
		}()
	}

	// Koodous
	if koodousAPI != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if kInfo, err := koodous.GetInfo(koodousAPI, app.Hashes.Sha256); err != nil {
				app.Errors.Koodous = err
			} else {
				app.Koodous = kInfo
			}
		}()
	}

	// VirusTotal - uses SHA256 hash for file lookup
	if vtAPIKey != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if vtInfo, err := vt.GetInfo(vtAPIKey, app.Hashes.Sha256); err != nil {
				app.Errors.VT = err
			} else {
				app.VirusTotal = vtInfo
			}
		}()
	}

	// Wait for all API calls to complete
	wg.Wait()

	return nil
}

// ExportJSON exports AndroidApp struct to a JSON file.
// It creates a pretty-printed JSON file with 2-space indentation.
// File permissions are set to 0600 (owner read/write only) for security.
func (app *AndroidApp) ExportJSON(jsonpath string) error {
	if jsonpath == "" {
		return fmt.Errorf("json export path cannot be empty")
	}

	dir := filepath.Dir(jsonpath)
	file, err := os.CreateTemp(dir, "."+filepath.Base(jsonpath)+".tmp-*")
	if err != nil {
		return fmt.Errorf("failed to create temporary JSON file for %q: %w", jsonpath, err)
	}
	tmpPath := file.Name()
	closed := false
	defer func() {
		if !closed {
			_ = file.Close()
		}
		_ = os.Remove(tmpPath)
	}()

	// Use json.Encoder for pretty-printed output with 2-space indent
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(app); err != nil {
		return fmt.Errorf("failed to encode JSON to %q: %w", jsonpath, err)
	}
	if err := file.Sync(); err != nil {
		return fmt.Errorf("failed to sync JSON file %q: %w", jsonpath, err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("failed to close JSON file %q: %w", jsonpath, err)
	}
	closed = true

	if err := os.Rename(tmpPath, jsonpath); err != nil {
		return fmt.Errorf("failed to replace JSON file %q: %w", jsonpath, err)
	}

	return nil
}
