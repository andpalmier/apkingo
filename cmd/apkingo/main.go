package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/shogo82148/androidbinary/apk"
)

/*
apkingo is a tool to dump detailed information about an apk file.
apkingo will explore the given file to get details on the apk, such as package
name, target SDK, permissions, metadata, certificate serial and issuer. The tool
will also retrieve information about the specified apk from the Play Store, and
(if valid API keys are provided) from Koodous and VirusTotal. In case the file
is not in VirusTotal, apkingo allows to upload it using the submitted API key.
*/

var err error

// main()
func main() {

	// check args
	if len(os.Args) != 2 {
		red.Println("[!] error, be sure to specify the file path: apkingo <apk_path>")
		os.Exit(1)
	}

	// read argument provided
	arg := os.Args[1]

	// print helper
	if arg == "-h" || arg == "help" || arg == "--help" {
		fmt.Println("Usage: apkingo <apk_path>")
		os.Exit(0)
	}

	// try to load apk from path
	pkg, err := apk.OpenFile(arg)
	if err != nil {
		red.Println("[!] error opening the apk file, be sure to specify the file path: apkingo <apk_path>")
		os.Exit(1)
	}
	defer pkg.Close()

	printBanner()

	// get API keys from env variables
	vtapi := os.Getenv("VT_API_KEY")
	kapi := os.Getenv("KOODOUS_API_KEY")
	if vtapi == "" {
		yellow.Println("[i] VirusTotal API key not found, you can export it using the env variable VT_API_KEY")
	}
	if kapi == "" {
		yellow.Println("[i] Koodous API key not found, you can export it using the env variable KOODOUS_API_KEY")
	}

	// extract general info
	generalInfo(*pkg)

	// extract hash values
	yellow.Println("\n* Hash values")
	sha256, err := hashInfo(arg)
	if err != nil {
		red.Printf("[!] error calculating hash values: %s\n", err.Error())
		os.Exit(1)
	}

	// extract permissions
	yellow.Println("\n* Permissions")
	permissionsInfo(*pkg)

	// extract metadata
	yellow.Println("\n* Metadata")
	metadataInfo(*pkg)

	// extract certificate info
	yellow.Println("\n* Certificate")
	if err := certInfo(arg); err != nil {
		red.Printf("[!] error while checking certificate: %s\n", err.Error())
	}

	// get Play Store info
	yellow.Println("\n* Play Store")
	if err := searchPlayStore(pkg.PackageName()); err != nil {
		red.Printf("%s\n", err)
	}

	// get Koodous info
	if kapi != "" {
		yellow.Println("\n* Koodous")
		err := koodousInfo(kapi, sha256)
		if err != nil {
			red.Printf("[!] error using Koodous - %s\n", err)
		}
	}

	// get VirusTotal info if api key is provided
	askvtupload := false
	if vtapi != "" {
		yellow.Println("\n* VirusTotal")
		err := vtInfo(vtapi, sha256)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				// ask to upload file to VT if it was not found
				askvtupload = true
				italic.Println("app not found in VirusTotal")
			} else {
				red.Printf("[!] error using VirusTotal - %s\n", err)
			}
		}
	}

	// ask to upload and scan file on VT
	if askvtupload {
		fmt.Printf("\n[i] do you want to upload the file to VT? (y/n): ")
		answer := bufio.NewScanner(os.Stdin)
		answer.Scan()
		answertext := strings.ToLower(answer.Text())
		if answertext == "yes" || answertext == "y" {
			vtScanFile(arg, vtapi)
		}
	}

	fmt.Println()
}
