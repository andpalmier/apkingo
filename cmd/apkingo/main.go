package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/andpalmier/apkingo/internal/analyzer"
	"github.com/andpalmier/apkingo/internal/config"
	"github.com/andpalmier/apkingo/internal/report"
	"github.com/andpalmier/apkingo/internal/ui"
	"github.com/andpalmier/apkingo/internal/vt"
)

func main() {
	cfg := config.Load()

	printer := ui.NewPrinter()
	reporter := report.NewReporter(printer)

	// Handle directory mode
	if cfg.DirPath != "" {
		processDirectory(cfg, printer, reporter)
		return
	}

	// Handle single file mode
	if cfg.APKPath == "" {
		log.Fatalf("No APK specified, please provide the path using the -apk flag.\n")
	}

	// Check if the input is an XAPK/APKS file
	if analyzer.IsXAPK(cfg.APKPath) {
		processXAPK(cfg, printer, reporter)
		return
	}

	// Regular APK processing
	app := analyzer.AndroidApp{}

	if err := app.ProcessAPK(cfg.APKPath, cfg.Country, cfg.VTAPIKey, cfg.KAPIKey, cfg.NoPlayStore); err != nil {
		log.Fatalf("Error processing APK: %v", err)
	}

	reporter.PrintBanner()
	reporter.PrintAll(&app)

	if cfg.JSONFile != "" {
		if err := app.ExportJSON(cfg.JSONFile); err != nil {
			log.Fatalf("Failed to export JSON: %v", err)
		}
	}

	if cfg.VTUpload {
		askToUploadToVT(cfg.APKPath, cfg.VTAPIKey)
	}
}

// processXAPK handles XAPK/APKS file processing.
// It extracts APKs from the archive, analyzes each one, and cleans up.
func processXAPK(cfg *config.Config, printer *ui.Printer, reporter *report.Reporter) {
	fmt.Printf("[i] Detected XAPK/APKS file: %s\n", cfg.APKPath)
	fmt.Println("[i] Extracting APKs from archive...")

	// Extract APKs from XAPK
	apkPaths, err := analyzer.ExtractAPKs(cfg.APKPath)
	if err != nil {
		log.Fatalf("Failed to extract APKs from XAPK: %v", err)
	}

	fmt.Printf("[i] Found %d APK(s) in archive\n\n", len(apkPaths))

	// Process each extracted APK with progress feedback (consistent with batch mode)
	var failed []string

	for i, apkPath := range apkPaths {
		fmt.Printf("[%d/%d] Processing: %s\n", i+1, len(apkPaths), filepath.Base(apkPath))

		app := analyzer.AndroidApp{}

		if err := app.ProcessAPK(apkPath, cfg.Country, cfg.VTAPIKey, cfg.KAPIKey, cfg.NoPlayStore); err != nil {
			fmt.Printf("  [!] Failed to process: %v\n", err)
			failed = append(failed, apkPath)
			continue
		}

		// Print detailed report for this APK
		reporter.PrintBanner()
		reporter.PrintAll(&app)

		// Export JSON if requested (consistent with batch mode)
		if cfg.JSONFile != "" {
			jsonPath := generateJSONPathForAPK(cfg.JSONFile, apkPath)
			if err := app.ExportJSON(jsonPath); err != nil {
				fmt.Printf("  [!] Failed to export JSON: %v\n", err)
			}
		}

		fmt.Printf("  [✓] Completed\n\n")
	}

	// Show failed count
	if len(failed) > 0 {
		fmt.Printf("[!] Failed to process %d out of %d APK(s)\n\n", len(failed), len(apkPaths))
	}

	// Clean up extracted files
	fmt.Println("[i] Cleaning up extracted files...")
	if err := analyzer.CleanupExtractedFiles(apkPaths); err != nil {
		log.Printf("Warning: Failed to clean up extracted files: %v", err)
	}

	fmt.Println("[i] Done!")
}

// processDirectory handles batch analysis of multiple APKs in a directory.
func processDirectory(cfg *config.Config, printer *ui.Printer, reporter *report.Reporter) {
	fmt.Printf("[i] Analyzing APKs in directory: %s\n", cfg.DirPath)

	results, failed, err := analyzer.AnalyzeDirectory(cfg.DirPath, cfg.Country, cfg.VTAPIKey, cfg.KAPIKey, cfg.NoPlayStore)
	if err != nil {
		log.Fatalf("Error analyzing directory: %v", err)
	}

	// Print detailed report for each APK (consistent with XAPK mode)
	if len(results) > 0 {
		fmt.Println("\n" + strings.Repeat("=", 60))
		fmt.Println("DETAILED ANALYSIS RESULTS")
		fmt.Println(strings.Repeat("=", 60))

		// Get ordered keys for consistent output
		keys := make([]string, 0, len(results))
		for k := range results {
			keys = append(keys, k)
		}

		for i, apkPath := range keys {
			app := results[apkPath]
			fmt.Printf("\n--- APK %d: %s ---\n", i+1, filepath.Base(apkPath))
			reporter.PrintAll(app)
		}
	}

	// Print summary report
	reporter.PrintBatchSummary(results, failed)

	// Export JSON if requested
	if cfg.JSONFile != "" {
		for apkPath, app := range results {
			jsonPath := generateJSONPathForAPK(cfg.JSONFile, apkPath)
			if err := app.ExportJSON(jsonPath); err != nil {
				log.Printf("Warning: Failed to export JSON for %s: %v", apkPath, err)
			}
		}
		fmt.Printf("\n[i] JSON reports exported to: %s\n", cfg.JSONFile)
	}

	fmt.Println("[i] Done!")
}

// generateJSONPathForAPK creates a JSON filename based on the APK filename.
func generateJSONPathForAPK(basePath, apkPath string) string {
	ext := ".json"
	base := strings.TrimSuffix(basePath, ext)
	// Get just the APK filename without extension
	apkName := strings.TrimSuffix(filepath.Base(apkPath), filepath.Ext(apkPath))
	return fmt.Sprintf("%s_%s%s", base, apkName, ext)
}

func askToUploadToVT(apkPath, vtAPIKey string) {
	fmt.Printf("\n[i] Do you want to upload the file to VirusTotal? (y/n): ")
	answer := bufio.NewScanner(os.Stdin)
	answer.Scan()
	ans := strings.ToLower(answer.Text())
	if ans == "yes" || ans == "y" {
		if err := vt.ScanFile(apkPath, vtAPIKey); err != nil {
			log.Fatalf("VirusTotal scan failed: %v", err)
		}
	}
}
