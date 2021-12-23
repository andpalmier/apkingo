package main

import (
	"fmt"

	"github.com/avast/apkverifier"
	"github.com/n0madic/google-play-scraper/pkg/app"
	"github.com/shogo82148/androidbinary/apk"
)

var androidname = map[int]string{
	1:  "Android 1",
	2:  "Android 1.1",
	3:  "Android 1.5",
	4:  "Android 1.6",
	5:  "Android 2",
	6:  "Android 2",
	7:  "Android 2.1",
	8:  "Android 2.2",
	9:  "Android 2.3",
	10: "Android 2.3.3",
	11: "Android 3",
	12: "Android 3.1",
	13: "Android 3.2",
	14: "Android 4",
	15: "Android 4.0.3",
	16: "Android 4.1",
	17: "Android 4.2",
	18: "Android 4.3",
	19: "Android 4.4",
	20: "Android 4.4W",
	21: "Android 5",
	22: "Android 5.1",
	23: "Android 6",
	24: "Android 7",
	25: "Android 7.1",
	26: "Android 8",
	27: "Android 8.1",
	28: "Android 9",
	29: "Android 10",
	30: "Android 11",
}

// AndroidApp -> includes apk.apk and app.app from the libraries
type AndroidApp struct {
	apk            apk.Apk
	playStoreInfo  app.App
	cert           apkverifier.CertInfo
	hashSHA256     []byte
	hashSHA1       []byte
	hashMD5        []byte
	path           string
	playStoreFound bool
}

// printAndroidInfo() - print all the information
func (androidapp *AndroidApp) printAndroidInfo() {
	// app name
	label, _ := androidapp.apk.Label(nil)
	if label != "" {
		fmt.Printf("\nApp name: %s\n", label)
	}
	fmt.Println("\n* General info")

	// package name
	fmt.Printf("Package name: %s\n", androidapp.apk.PackageName())

	// app version
	version, _ := androidapp.apk.Manifest().VersionName.String()
	if version != "" {
		fmt.Printf("Version Name: %s\n", version)
	}

	// Main activity
	mainactivity, _ := androidapp.apk.MainActivity()
	if mainactivity != "" {
		fmt.Printf("Main Activity name: %s\n", mainactivity)
	}

	// Target and Minimum SDK
	sdktarget, _ := androidapp.apk.Manifest().SDK.Target.Int32()
	if string(sdktarget) != "" {
		fmt.Printf("Target SDK: %d (%s)\n", sdktarget, androidname[int(sdktarget)])
	}
	sdkmin, _ := androidapp.apk.Manifest().SDK.Min.Int32()
	if string(sdkmin) != "" {
		fmt.Printf("Minimum SDK: %d (%s)\n", sdkmin, androidname[int(sdkmin)])
	}

	// Hash values
	fmt.Println("\n* Hash values")
	fmt.Printf("MD5: %x\n", androidapp.hashMD5)
	fmt.Printf("SHA1: %x\n", androidapp.hashSHA1)
	fmt.Printf("SHA256: %x\n", androidapp.hashSHA256)

	// Permissions
	fmt.Println("\n* Permissions")
	for _, n := range androidapp.apk.Manifest().UsesPermissions {
		permission, _ := n.Name.String()
		if permission != "" {
			fmt.Printf("%s\n", permission)
		}
	}

	// Metadata
	fmt.Println("\n* Metadata")
	if len(androidapp.apk.Manifest().App.MetaData) > 0 {
		for _, n := range androidapp.apk.Manifest().App.MetaData {
			metaname, _ := n.Name.String()
			metavalue, _ := n.Value.String()
			if metaname != "" {
				if metavalue != "" {
					fmt.Printf("%s: %s\n", metaname, metavalue)
				} else {
					fmt.Printf("%s\n", metaname)
				}
			}
		}
	} else {
		fmt.Println("no metadata found")
	}

	// Certificate info
	if cert {
		fmt.Println("\n* Certificate info:")
		if androidapp.cert.Issuer != "" {
			fmt.Printf("SHA1: %s\n", androidapp.cert.Sha1)
			fmt.Printf("Issuer: %s\n", androidapp.cert.Issuer)
			fmt.Printf("Valid from: %s\n", androidapp.cert.ValidFrom)
			fmt.Printf("Valid to: %s\n", androidapp.cert.ValidTo)
		} else {
			fmt.Println("Certificate not found")
		}
	}

	// Play Store info
	if playstore {
		fmt.Println("\n* Play Store info")
		if androidapp.playStoreFound {
			fmt.Printf("URL: %s\n", androidapp.playStoreInfo.URL)
			fmt.Printf("Version: %s\n", androidapp.playStoreInfo.Version)
			fmt.Printf("App summary: %s\n", androidapp.playStoreInfo.Summary)
			fmt.Printf("App developer: %s (email: %s, id: %s)\n", androidapp.playStoreInfo.Developer, androidapp.playStoreInfo.DeveloperEmail, androidapp.playStoreInfo.DeveloperID)
			fmt.Printf("Release date: %s\n", androidapp.playStoreInfo.Released)
			fmt.Printf("Installs: %s\n", androidapp.playStoreInfo.Installs)
			fmt.Printf("Score: %f\n", androidapp.playStoreInfo.Score)
		} else {
			fmt.Println("app not found in Play Store")
		}
	}

	fmt.Println()
}
