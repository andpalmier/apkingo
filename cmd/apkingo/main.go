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
	vtUpload = false
)

const (
	kAPImsg = "[i] Koodous API key not found, you can provide it with the -kapi flag or " +
		"export it using the env variable KOODOUS_API_KEY"
	vtAPImsg = "[i] VirusTotal API key not found, you can provide it with the -vtapi flag or" +
		"export it using the env variable VT_API_KEY)"
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
		log.Fatalf("No APK specified, please provide the path using the -apk flag.\n")
	}

	vtAPIKey = getAPIKey(vtAPIKey, "VT_API_KEY", vtAPImsg)
	kAPIKey = getAPIKey(kAPIKey, "KOODOUS_API_KEY", kAPImsg)

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
		return fmt.Errorf("error loading APK: %s", err)
	}
	defer pkg.Close()

	if err = app.setGeneralInfo(pkg); err != nil {
		log.Printf("error getting general information: %s\n", err)
	}

	if err = app.setHashes(apkPath); err != nil {
		return fmt.Errorf("error setting hashes: %s\n", err)
	}

	if err = app.setCertInfo(apkPath); err != nil {
		log.Printf("error getting certificate information: %s\n", err)
	}

	if err = app.setPlayStoreInfo(country); err != nil {
		log.Printf("error getting Play Store information: %s\n", err)
	}

	if koodousAPI != "" {
		if err = app.setKoodousInfo(koodousAPI); err != nil {
			log.Printf("error getting Koodous information: %s\n", err)
		}
	}

	if vtAPIKey != "" {
		if err = app.setVTInfo(vtAPIKey); err != nil {
			log.Printf("error getting VirusTotal information: %s\n", err)
			vtUpload = true
		}
	}

	return nil
}

func loadAPK(filePath string) (*apk.Apk, error) {
	return apk.OpenFile(filePath)
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
