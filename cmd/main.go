package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/shogo82148/androidbinary/apk"
)

/*

apkingo is a tool written in Go to get detailed information about an apk file. apkingo will
explore the given file to get details on the apk, such as package name, target SDK, permissions
and metadata.

After downloading the repository, navigate into the directory and build the project with make
apkingo. This will create a folder build, containing an executable called apkingo. You can then
run the executable with the following flags:

-apk to specify the path to the apk file (required)
-cert for printing the certificate information retrieved in the apk file (sometimes it returns
	a conversion error, but it's still working!)
-playstore for searching the app in the Play Store by its package name

*/

var err error
var androidapp AndroidApp
var apkpath string
var cert bool
var playstore bool

func init() {
	flag.StringVar(&apkpath, "apk", "", "specify apk path")
	flag.BoolVar(&cert, "cert", false, "get certificate info")
	flag.BoolVar(&playstore, "playstore", false, "search app in Play Store by package name")
}

func main() {
	flag.Parse()

	// load the apk
	androidapp.path = apkpath
	pkg, err := apk.OpenFile(androidapp.path)
	if err != nil {
		fmt.Println("error opening the apk file, be sure to use '-apk' to specify the file path")
		os.Exit(1)
	}
	defer pkg.Close()
	androidapp.apk = *pkg

	// extract hash values
	if err = androidapp.getHashValues(); err != nil {
		fmt.Printf("error calculating hash values: %s\n", err.Error())
		os.Exit(1)
	}

	// get certificate info
	if cert {
		if err = androidapp.getCertInfo(); err != nil {
			fmt.Printf("error retrieving certificate information: %s\n", err.Error())
		}
	}

	// search the app in the Play Store by package name
	if playstore {
		if err = androidapp.searchPlayStore(); err != nil {
			androidapp.playStoreFound = false
		}
	}

	// print result
	androidapp.printAndroidInfo()
}
