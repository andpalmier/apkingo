package analyzer

import (
	"archive/zip"
	"sort"
	"strings"

	"github.com/andpalmier/apkingo/internal/utils"
	"github.com/shogo82148/androidbinary"
	"github.com/shogo82148/androidbinary/apk"
)

// canonicalArchOrder defines the preferred display order for Android ABIs.
// Pure-Java apps will have an empty list.
var canonicalArchOrder = map[string]int{
	"arm64-v8a":   0,
	"armeabi-v7a": 1,
	"armeabi":     2,
	"x86_64":      3,
	"x86":         4,
}

// SetGeneralInfo sets general information about the APK.
// If locale is non-empty, it also extracts the app name for that locale
// (e.g., "en", "zh-CN") and stores it in NameLocale.
func (app *AndroidApp) SetGeneralInfo(pkg *apk.Apk, locale string) error {
	var err error

	app.Name, err = pkg.Label(nil)
	if err != nil {
		return err
	}

	// Extract locale-specific app name if requested
	if locale != "" {
		resConfig := androidbinaryParseLocale(locale)
		if localizedName, err := pkg.Label(resConfig); err == nil && localizedName != "" && localizedName != app.Name {
			app.NameLocale = localizedName
			app.Locale = locale
		} else if err != nil {
			utils.LogError("error getting localized app name", err)
		}
	}

	app.PackageName = pkg.PackageName()

	// Extract version name (human-readable, e.g. "1.2.3")
	app.VersionName, err = pkg.Manifest().VersionName.String()
	utils.LogError("error getting version name", err)

	// Extract version code (numeric build identifier, e.g. 42)
	app.VersionCode, err = pkg.Manifest().VersionCode.Int32()
	utils.LogError("error getting version code", err)

	app.MainActivity, err = pkg.MainActivity()
	utils.LogError("error getting main activity information", err)

	app.MinimumSDK, err = pkg.Manifest().SDK.Min.Int32()
	utils.LogError("error getting minimum SDK information", err)

	app.TargetSDK, err = pkg.Manifest().SDK.Target.Int32()
	utils.LogError("error getting target SDK information", err)

	for _, n := range pkg.Manifest().UsesPermissions {
		permission, _ := n.Name.String()
		if permission != "" {
			app.Permissions = append(app.Permissions, permission)
		}
	}

	for _, n := range pkg.Manifest().App.MetaData {
		metadataName, _ := n.Name.String()
		metadataValue, _ := n.Value.String()
		if metadataName != "" {
			app.Metadata = append(app.Metadata, Metadata{
				Name:  metadataName,
				Value: metadataValue,
			})
		}
	}

	return nil
}

// SetArchitectures detects supported CPU architectures by scanning the APK's
// lib/ directory structure (e.g., arm64-v8a, armeabi-v7a, x86_64).
// Pure-Java/Kotlin apps with no native code will have an empty slice.
func (app *AndroidApp) SetArchitectures(apkPath string) {
	r, err := zip.OpenReader(apkPath)
	if err != nil {
		utils.LogError("error scanning APK for native architectures", err)
		return
	}
	defer r.Close()

	app.Architectures = make([]string, 0)
	seen := make(map[string]bool)
	for _, f := range r.File {
		// Look for entries under lib/<arch>/ — the standard Android ABI layout
		if strings.HasPrefix(f.Name, "lib/") {
			parts := strings.SplitN(f.Name, "/", 3)
			if len(parts) >= 2 && parts[1] != "" && !seen[parts[1]] {
				app.Architectures = append(app.Architectures, parts[1])
				seen[parts[1]] = true
			}
		}
	}

	// Sort in canonical ABI order for consistent output
	sort.Slice(app.Architectures, func(i, j int) bool {
		oi, oki := canonicalArchOrder[app.Architectures[i]]
		oj, okj := canonicalArchOrder[app.Architectures[j]]
		if !oki {
			oi = len(canonicalArchOrder)
		}
		if !okj {
			oj = len(canonicalArchOrder)
		}
		return oi < oj
	})
}

// androidbinaryParseLocale converts a locale string (e.g., "en", "zh-CN", "en-US")
// into an androidbinary.ResTableConfig suitable for localized resource lookups.
//
// The format follows Android's resource configuration locale format:
//   - Language: ISO 639-1 two-letter code (lowercase), required
//   - Country:  ISO 3166-1-alpha-2 two-letter code (uppercase), optional
//
// If the locale string is empty or malformed, nil is returned (default locale).
func androidbinaryParseLocale(locale string) *androidbinary.ResTableConfig {
	config := &androidbinary.ResTableConfig{}
	parts := strings.SplitN(locale, "-", 2)

	// Language part: first 2 chars, lowercase
	lang := strings.ToLower(parts[0])
	if len(lang) >= 2 {
		config.Language = [2]uint8{lang[0], lang[1]}
	}

	// Country part: if present, next 2 chars, uppercase
	if len(parts) > 1 {
		country := strings.ToUpper(parts[1])
		if len(country) >= 2 {
			config.Country = [2]uint8{country[0], country[1]}
		}
	}

	return config
}
