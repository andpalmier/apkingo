package config

import (
	"flag"
	"fmt"
	"os"
)

const (
	kAPImsg = "[i] Koodous API key not found, you can provide it with the -kapi flag or " +
		"export it using the env variable KOODOUS_API_KEY"
	vtAPImsg = "[i] VirusTotal API key not found, you can provide it with the -vtapi flag or " +
		"export it using the env variable VT_API_KEY)"
)

type Config struct {
	APKPath  string
	JSONFile string
	Country  string
	VTAPIKey string
	KAPIKey  string
	VTUpload bool
	NoColor  bool
}

func Load() *Config {
	cfg := &Config{}
	flag.StringVar(&cfg.JSONFile, "json", "", "Path to export analysis in JSON format")
	flag.StringVar(&cfg.APKPath, "apk", "", "Path to APK file")
	flag.StringVar(&cfg.Country, "country", "us", "Country code of the Play Store")
	flag.StringVar(&cfg.VTAPIKey, "vtapi", "", "VirusTotal API key")
	flag.StringVar(&cfg.KAPIKey, "kapi", "", "Koodous API key")
	flag.BoolVar(&cfg.NoColor, "nocolor", false, "Disable colored output")
	flag.Parse()

	if len(cfg.Country) != 2 {
		// Default to "us" if invalid
		cfg.Country = "us"
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
