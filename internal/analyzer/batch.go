// Package analyzer provides APK analysis functionality.
package analyzer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/andpalmier/apkingo/v3/internal/utils"
)

// AnalyzeDirectory analyzes all APK files in a given directory.
// It takes country, locale, and API keys for external service lookups (Play Store, VirusTotal, Koodous).
// noPlayStore when true, skips Play Store API calls for offline analysis.
// Returns a map of file paths to AndroidApp results and a slice of failed paths.
func AnalyzeDirectory(dirPath, country, locale, vtAPIKey, koodousAPI string, noPlayStore bool) (map[string]*AndroidApp, []string, error) {
	// Check if directory exists
	info, err := os.Stat(dirPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to access directory: %w", err)
	}

	if !info.IsDir() {
		return nil, nil, fmt.Errorf("path is not a directory: %s", dirPath)
	}

	// Find all APK files in the directory
	apkPaths, err := findAPKs(dirPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to find APKs: %w", err)
	}

	if len(apkPaths) == 0 {
		return nil, nil, fmt.Errorf("no APK files found in directory: %s", dirPath)
	}

	results := make(map[string]*AndroidApp)
	var failed []string
	var analysisPaths []string
	var extractedPaths []string

	for _, apkPath := range apkPaths {
		if !IsXAPK(apkPath) {
			analysisPaths = append(analysisPaths, apkPath)
			continue
		}

		paths, err := ExtractAPKs(apkPath)
		if err != nil {
			failed = append(failed, apkPath)
			fmt.Printf("  [!] Failed to extract %s: %s\n", utils.SanitizeText(filepath.Base(apkPath)), utils.SanitizeText(err.Error()))
			continue
		}
		analysisPaths = append(analysisPaths, paths...)
		extractedPaths = append(extractedPaths, paths...)
	}

	if len(extractedPaths) > 0 {
		defer func() {
			if err := CleanupExtractedFiles(extractedPaths); err != nil {
				fmt.Printf("[!] Failed to clean up extracted files: %s\n", utils.SanitizeText(err.Error()))
			}
		}()
	}

	// Analyze each APK with progress feedback
	fmt.Printf("[i] Found %d APK(s) to analyze\n\n", len(analysisPaths))

	for i, apkPath := range analysisPaths {
		fmt.Printf("[%d/%d] Processing: %s\n", i+1, len(analysisPaths), utils.SanitizeText(filepath.Base(apkPath)))

		app := &AndroidApp{}
		if err := app.ProcessAPK(apkPath, country, locale, vtAPIKey, koodousAPI, noPlayStore); err != nil {
			// Log error but continue with other files
			fmt.Printf("  [!] Failed to process: %s\n", utils.SanitizeText(err.Error()))
			failed = append(failed, apkPath)
			continue
		}
		results[apkPath] = app
		fmt.Printf("  [✓] Completed\n\n")
	}

	if len(failed) > 0 {
		fmt.Printf("[!] Failed to process %d out of %d APK(s)\n\n", len(failed), len(analysisPaths))
	}

	return results, failed, nil
}

// findAPKs returns a list of all APK file paths in the given directory.
func findAPKs(dirPath string) ([]string, error) {
	var apkPaths []string

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Check file extension
		name := entry.Name()
		switch strings.ToLower(filepath.Ext(name)) {
		case ".apk", ".xapk", ".apks":
			apkPaths = append(apkPaths, filepath.Join(dirPath, name))
		}
	}

	return apkPaths, nil
}
