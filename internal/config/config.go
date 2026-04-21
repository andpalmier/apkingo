package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/andpalmier/apkingo/internal/constants"
)

const (
	kAPImsg = "[i] Koodous API key not found, you can provide it with the -kapi flag or " +
		"export it using the env variable KOODOUS_API_KEY"
	vtAPImsg = "[i] VirusTotal API key not found, you can provide it with the -vtapi flag or " +
		"export it using the env variable VT_API_KEY)"
)

type Config struct {
	APKPath     string
	DirPath     string
	JSONFile    string
	Country     string
	VTAPIKey    string
	KAPIKey     string
	VTUpload    bool
	NoPlayStore bool
}

func Load() *Config {
	cfg := &Config{}
	flag.StringVar(&cfg.JSONFile, "json", "", "Path to export analysis in JSON format")
	flag.StringVar(&cfg.APKPath, "apk", "", "Path to APK or XAPK file")
	flag.StringVar(&cfg.DirPath, "dir", "", "Analyze all APKs in a directory")
	flag.StringVar(&cfg.Country, "country", constants.DefaultCountry, "Country code of the Play Store")
	flag.StringVar(&cfg.VTAPIKey, "vtapi", "", "VirusTotal API key")
	flag.StringVar(&cfg.KAPIKey, "kapi", "", "Koodous API key")
	flag.BoolVar(&cfg.VTUpload, "vtupload", false, "Upload APK to VirusTotal after analysis")
	flag.BoolVar(&cfg.NoPlayStore, "no-play-store", false, "Skip Play Store API calls for offline analysis")
	flag.Parse()

	// Validate input options
	if cfg.APKPath != "" && cfg.DirPath != "" {
		fmt.Println("Error: Cannot specify both -apk and -dir flags")
		fmt.Println("Please use either -apk for a single file or -dir for a directory")
		os.Exit(1)
	}

	if cfg.APKPath == "" && cfg.DirPath == "" {
		fmt.Println("Error: No APK or directory specified")
		fmt.Println("Please provide either -apk <path> or -dir <path>")
		os.Exit(1)
	}

	if len(cfg.Country) != 2 {
		// Default to constants.DefaultCountry if invalid
		cfg.Country = constants.DefaultCountry
	}

	cfg.VTAPIKey = getAPIKey(cfg.VTAPIKey, "VT_API_KEY", vtAPImsg)
	cfg.KAPIKey = getAPIKey(cfg.KAPIKey, "KOODOUS_API_KEY", kAPImsg)

	return cfg
}

func getAPIKey(flagValue, envVar, msg string) string {
	if flagValue == "" {
		flagValue = os.Getenv(envVar)
	}
	if flagValue == "" && msg != "" {
		fmt.Println(msg)
	}
	return flagValue
}
