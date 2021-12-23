package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/shogo82148/androidbinary/apk"
)

// TODO SDK Name?

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
	androidapp.Path = apkpath
	pkg, err := apk.OpenFile(androidapp.Path)
	if err != nil {
		fmt.Println("error opening the apk file, be sure to use '-apk' to specify the file path")
		os.Exit(1)
	}
	defer pkg.Close()
	androidapp.Apk = *pkg

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
			androidapp.PlayStoreFound = false
		}
	}

	// print result
	androidapp.printAndroidInfo()
}
