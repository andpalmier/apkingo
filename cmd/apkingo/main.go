package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/shogo82148/androidbinary/apk"
)

var (
	apkPath  string
	jsonFile string
	country  string
	vtAPIKey string
	kAPIKey  string
	kAPImsg  = "[i] Koodous API key not found, you can provide it with the -kapi flag or " +
		"export it using the env variable KOODOUS_API_KEY"
	vtAPImsg = "[i] VirusTotal API key not found, you can provide it with the -vtapi flag or" +
		"export it using the env variable VT_API_KEY)"
	vtUpload = false
)

func init() {
	flag.StringVar(&jsonFile, "json", "", "Path to export analysis in JSON format")
	flag.StringVar(&apkPath, "apk", "", "Path to APK file")
	flag.StringVar(&country, "country", "us", "Country code of the Play Store")
	flag.StringVar(&vtAPIKey, "vtapi", "", "VirusTotal API key")
	flag.StringVar(&kAPIKey, "kapi", "", "Koodous API key")

}

func main() {
	flag.Parse()
	if apkPath == "" {
		fmt.Println("No APK specified, please provide the path using the -apk flag.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if vtAPIKey == "" {
		vtAPIKey = getAPIKey("VT_API_KEY", vtAPImsg)
	}
	if kAPIKey == "" {
		kAPIKey = getAPIKey("KOODOUS_API_KEY", kAPImsg)
	}
	app := AndroidApp{}
	if err := app.processAPK(apkPath, country, vtAPIKey, kAPIKey); err != nil {
		log.Fatalf("Error processing APK: %v", err)
	}

	printBanner()
	printAll(&app)

	if jsonFile != "" {
		app.ExportJSON(jsonFile)
	}

	if vtUpload {
		askToUploadToVT(apkPath, vtAPIKey)
	}
}

func (app *AndroidApp) processAPK(apkPath, country, vtAPIKey, koodousAPI string) error {
	pkg, err := loadAPK(apkPath)
	if err != nil {
		return fmt.Errorf("error loading APK: %w", err)
	}
	defer pkg.Close()

	app.setGeneralInfo(pkg)

	if err = app.setHashes(apkPath); err != nil {
		return fmt.Errorf("error setting hashes: %w", err)
	}

	if err = app.setCertInfo(apkPath); err != nil {
		log.Printf("error with certificate: %w", err)
	}

	app.setPlayStoreInfo(country)

	if koodousAPI != "" {
		app.setKoodousInfo(koodousAPI)
	}

	if vtAPIKey != "" {
		if err := app.setVTInfo(vtAPIKey); err != nil {
			vtUpload = true
		}
	}

	return nil
}

func loadAPK(filePath string) (*apk.Apk, error) {
	pkg, err := apk.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	return pkg, nil
}

func askToUploadToVT(apkPath, vtAPIKey string) {
	fmt.Printf("\n[i] Do you want to upload the file to VirusTotal? (y/n): ")
	answer := bufio.NewScanner(os.Stdin)
	answer.Scan()
	ans := strings.ToLower(answer.Text())
	if ans == "yes" || ans == "y" {
		vtScanFile(apkPath, vtAPIKey)
	}
}
