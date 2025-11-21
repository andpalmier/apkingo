package analyzer

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/andpalmier/apkingo/internal/koodous"
	"github.com/andpalmier/apkingo/internal/vt"
	"github.com/shogo82148/androidbinary/apk"
)

// AndroidApp represents information extracted from an APK file
type AndroidApp struct {
	Name         string               `json:"name"`
	PackageName  string               `json:"package-name"`
	Version      string               `json:"version"`
	MainActivity string               `json:"main-activity"`
	MinimumSDK   int32                `json:"minimum-sdk"`
	TargetSDK    int32                `json:"target-sdk"`
	Hashes       Hashes               `json:"hashes"`
	Permissions  []string             `json:"permissions"`
	Metadata     []Metadata           `json:"metadata"`
	Certificate  CertificateInfo      `json:"certificate"`
	PlayStore    *PlayStoreInfo       `json:"playstore,omitempty"`
	Koodous      *koodous.KoodousInfo `json:"koodous,omitempty"`
	VirusTotal   *vt.VirusTotalInfo   `json:"virustotal,omitempty"`
	Errors       AnalysisErrors       `json:"-"` // internal use only
}

// AnalysisErrors holds errors encountered during analysis
type AnalysisErrors struct {
	General   error
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

// ProcessAPK orchestrates the APK analysis
func (app *AndroidApp) ProcessAPK(apkPath, country, vtAPIKey, koodousAPI string) error {
	pkg, err := apk.OpenFile(apkPath)
	if err != nil {
		return fmt.Errorf("error loading APK: %s", err)
	}
	defer pkg.Close()

	if err = app.SetGeneralInfo(pkg); err != nil {
		app.Errors.General = err
	}

	if err = app.SetHashes(apkPath); err != nil {
		return fmt.Errorf("error setting hashes: %s", err)
	}

	if err = app.SetCertInfo(apkPath); err != nil {
		app.Errors.Cert = err
	}

	// Run external API calls concurrently for better performance
	var wg sync.WaitGroup

	// PlayStore
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := app.SetPlayStoreInfo(country); err != nil {
			app.Errors.PlayStore = err
		}
	}()

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

	// VirusTotal
	if vtAPIKey != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if vtInfo, err := vt.GetInfo(vtAPIKey, app.Hashes.Sha1); err != nil {
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

// ExportJSON exports AndroidApp struct to a JSON file
func (app *AndroidApp) ExportJSON(jsonpath string) error {
	// Validate path
	if jsonpath == "" {
		return fmt.Errorf("json export path cannot be empty")
	}

	// Create file with secure permissions (0600 - owner read/write only)
	file, err := os.OpenFile(jsonpath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to create JSON file %q: %w", jsonpath, err)
	}
	defer file.Close()

	// Use json.Encoder for more efficient streaming
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")

	if err := encoder.Encode(app); err != nil {
		return fmt.Errorf("failed to encode JSON to %q: %w", jsonpath, err)
	}

	return nil
}
