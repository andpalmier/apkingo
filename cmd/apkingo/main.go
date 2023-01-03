package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/shogo82148/androidbinary/apk"
)

/*
apkingo is a tool written in Go to get detailed information about an apk file.
apkingo will explore the given file to get details on the apk, such as package
name, target SDK, permissions, metadata, certificate serial and issuer.
The tool will also retrieve information about the specified apk from the
Play Store, and (if valid API keys are provided) from Koodous and VirusTotal.

## Installation

You can can download apkingo from the Releases section or compile it from the
source by downloading the repository, navigating into the "apkingo" directory
and building the project with "make apkingo". This will create a "build" folder,
containing the resulting executable.

## Usage

You can run apkingo with the following flags:

- -apk to specify the path to the apk file (**required**)
- -json to specify the path of the json file where the results will be exported
- -vt to specify VirusTotal API key (required for VirusTotal analysis)
- -k to specify Koodous API key (required for Koodous analysis)

*/

var err error
var androidapp AndroidApp
var apkpath, jsonfile, vtapi, kapi string

// init() - parse flags
func init() {
	flag.StringVar(&apkpath, "apk", "", "specify apk path")
	flag.StringVar(&jsonfile, "json", "", "specify json file path to export findings")
	flag.StringVar(&vtapi, "vt", "", "specify VirusTotal API key (required for VT analysis)")
	flag.StringVar(&kapi, "k", "", "specify Koodous API key (required for Koodous analysis)")
}

// main()
func main() {
	printBanner()
	flag.Parse()

	// load the apk
	pkg, err := apk.OpenFile(apkpath)
	if err != nil {
		red.Println("[!] error opening the apk file, be sure to use '-apk' to specify the file path")
		os.Exit(1)
	}
	defer pkg.Close()

	androidapp.setApkGeneralInfo(*pkg)
	androidapp.setPermissions(*pkg)
	androidapp.setMetadata(*pkg)

	// extract hash values
	if err = androidapp.setHashValues(apkpath); err != nil {
		red.Printf("[!] error calculating hash values: %s\n", err.Error())
		os.Exit(1)
	}

	// extract certificate info
	certinfo, err := androidapp.getCertInfo(apkpath)
	if err != nil {
		androidapp.Certificate.Issuer = ""
	} else {
		androidapp.setCertInfo(*certinfo)
	}

	// get Play Store info
	if playstoreinfo, err := androidapp.searchPlayStore(); err != nil {
		androidapp.PlayStore = nil
	} else {
		androidapp.setPlayStoreInfo(playstoreinfo)
	}

	// get Koodous info
	if kapi != "" {
		err = androidapp.getKoodousDetection(kapi)
		if err != nil {
			red.Printf("[!] error with Koodous API: %s\n", err)
			kapi = ""
		}
	}

	// get VirusTotal info if api key is provided
	if vtapi != "" {
		err = androidapp.getVTDetection(vtapi)
		if err != nil {
			red.Printf("[!] error with VirusTotal API: %s\n", err)
			vtapi = ""
		}
	}

	// print result
	androidapp.printAll()

	// save result as json if json path is provided
	if jsonfile != "" {
		if filepath.Ext(jsonfile) != ".json" {
			jsonfile = jsonfile + ".json"
		}
		err = androidapp.ExportJson(jsonfile)
		if err == nil {
			fmt.Printf("analysis exported at %s\n", jsonfile)
		} else {
			red.Printf("[!] error exporting results at %s: %s", jsonfile, err)
		}
	}
}
