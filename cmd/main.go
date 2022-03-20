package main

import (
	"flag"
	"fmt"
	"github.com/shogo82148/androidbinary/apk"
	"os"
	"path/filepath"
)

/*

apkingo is a tool written in Go to get detailed information about an apk file. apkingo will
explore the given file to get details on the apk, such as package name, target SDK, permissions,
metadata, certificate serial and issuer. The tool will also retrieve information about the
specified apk from the Play Store and Koodous.

After downloading the repository, navigate into the directory and build the project with make
apkingo. This will create a folder build, containing an executable called apkingo. You can then
run the executable with the following flags:

-apk 	to specify the path to the apk file (required)
-json	to specify the path of the json file where the results will be exported (optional)
-vt		to specify VirusTotal api (optional)

*/

var err error
var androidapp AndroidApp
var apkpath string
var jsonfile string
var vtapi string

// init() - parse flags
func init() {
	flag.StringVar(&apkpath, "apk", "", "specify apk path")
	flag.StringVar(&jsonfile, "json", "", "specify json file to export findings in json")
	flag.StringVar(&vtapi, "vt", "", "specify VirusTotal api (required for VT analysis)")
}

// main()
func main() {
	flag.Parse()

	// load the apk
	pkg, err := apk.OpenFile(apkpath)
	if err != nil {
		fmt.Println("error opening the apk file, be sure to use '-apk' to specify the file path")
		os.Exit(1)
	}
	defer pkg.Close()
	androidapp.setApkGeneralInfo(*pkg)
	androidapp.setPermissions(*pkg)
	androidapp.setMetadata(*pkg)

	// extract hash values
	if err = androidapp.setHashValues(apkpath); err != nil {
		fmt.Printf("error calculating hash values: %s\n", err.Error())
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
		androidapp.PlayStore.Url = ""
	} else {
		androidapp.setPlayStoreInfo(playstoreinfo)
	}
	// get Koodous info
	err = androidapp.getKoodousDetection()
	if err != nil {
		androidapp.Koodous.Url = ""
		fmt.Printf("%s", err)
	}

	if vtapi != "" {
		err = androidapp.getVTDetection(vtapi)
		if err != nil {
			fmt.Println(err)
		}
	}

	// print result
	androidapp.printAll()

	// save result as json
	if jsonfile != "" {
		if filepath.Ext(jsonfile) != ".json" {
			jsonfile = jsonfile + ".json"
		}
		androidapp.ExportJson(jsonfile)
	}
}
