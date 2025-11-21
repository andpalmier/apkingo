package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/andpalmier/apkingo/internal/analyzer"
	"github.com/andpalmier/apkingo/internal/config"
	"github.com/andpalmier/apkingo/internal/report"
	"github.com/andpalmier/apkingo/internal/ui"
	"github.com/andpalmier/apkingo/internal/vt"
)

func main() {
	cfg := config.Load()

	if cfg.APKPath == "" {
		log.Fatalf("No APK specified, please provide the path using the -apk flag.\n")
	}

	printer := ui.NewPrinter(cfg.NoColor)
	reporter := report.NewReporter(printer)

	app := analyzer.AndroidApp{}

	if err := app.ProcessAPK(cfg.APKPath, cfg.Country, cfg.VTAPIKey, cfg.KAPIKey); err != nil {
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
