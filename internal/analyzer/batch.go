// Package analyzer provides APK analysis functionality.
package analyzer

import (
	"fmt"
	"os"
	"path/filepath"
)

// AnalyzeDirectory analyzes all APK files in a given directory.
// It takes country and API keys for external service lookups (VirusTotal, Koodous).
// Returns a map of file paths to AndroidApp results and a slice of failed paths.
func AnalyzeDirectory(dirPath, country, vtAPIKey, koodousAPI string) (map[string]*AndroidApp, []string, error) {
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

	// Analyze each APK with progress feedback
	results := make(map[string]*AndroidApp)
	var failed []string

	fmt.Printf("[i] Found %d APK(s) to analyze\n\n", len(apkPaths))

	for i, apkPath := range apkPaths {
		fmt.Printf("[%d/%d] Processing: %s\n", i+1, len(apkPaths), filepath.Base(apkPath))

		app := &AndroidApp{}
		if err := app.ProcessAPK(apkPath, country, vtAPIKey, koodousAPI); err != nil {
			// Log error but continue with other files
			fmt.Printf("  [!] Failed to process: %v\n", err)
			failed = append(failed, apkPath)
			continue
		}
		results[apkPath] = app
		fmt.Printf("  [✓] Completed\n\n")
	}

	if len(failed) > 0 {
		fmt.Printf("[!] Failed to process %d out of %d APK(s)\n\n", len(failed), len(apkPaths))
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
		ext := filepath.Ext(name)
		if ext == ".apk" || ext == ".xapk" || ext == ".apks" {
			apkPaths = append(apkPaths, filepath.Join(dirPath, name))
		}
	}

	return apkPaths, nil
}
