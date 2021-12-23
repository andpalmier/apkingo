package main

import (
	"fmt"

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

type AndroidApp struct {
	Apk            apk.Apk
	PlayStoreInfo  app.App
	Cert           string
	HashSHA256     []byte
	HashSHA1       []byte
	HashMD5        []byte
	Path           string
	PlayStoreFound bool
}

// printAndroidInfo() - print all the information
func (androidapp *AndroidApp) printAndroidInfo() {
	// app name
	label, _ := androidapp.Apk.Label(nil)
	if label != "" {
		fmt.Printf("\nApp name: %s\n", label)
	}
	fmt.Println("\n* General info")

	// package name
	fmt.Printf("Package name: %s\n", androidapp.Apk.PackageName())

	// app version
	version, _ := androidapp.Apk.Manifest().VersionName.String()
	if version != "" {
		fmt.Printf("Version Name: %s\n", version)
	}

	// Main activity
	mainactivity, _ := androidapp.Apk.MainActivity()
	if mainactivity != "" {
		fmt.Printf("Main Activity name: %s\n", mainactivity)
	}

	// Target and Minimum SDK
	sdktarget, _ := androidapp.Apk.Manifest().SDK.Target.Int32()
	if string(sdktarget) != "" {
		fmt.Printf("Target SDK: %d (%s)\n", sdktarget, androidname[int(sdktarget)])
	}
	sdkmin, _ := androidapp.Apk.Manifest().SDK.Min.Int32()
	if string(sdkmin) != "" {
		fmt.Printf("Minimum SDK: %d (%s)\n", sdkmin, androidname[int(sdkmin)])
	}

	// Hash values
	fmt.Println("\n* Hash values")
	fmt.Printf("MD5: %x\n", androidapp.HashMD5)
	fmt.Printf("SHA1: %x\n", androidapp.HashSHA1)
	fmt.Printf("SHA256: %x\n", androidapp.HashSHA256)

	// Permissions
	fmt.Println("\n* Permissions")
	for _, n := range androidapp.Apk.Manifest().UsesPermissions {
		permission, _ := n.Name.String()
		if permission != "" {
			fmt.Printf("%s\n", permission)
		}
	}

	// Metadata
	fmt.Println("\n* Metadata")
	if len(androidapp.Apk.Manifest().App.MetaData) > 0 {
		for _, n := range androidapp.Apk.Manifest().App.MetaData {
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
		fmt.Printf("\n* Certificate info:\n%s\n", androidapp.Cert)
	}

	// Play Store info
	if playstore {
		fmt.Println("\n* Play Store info")
		if androidapp.PlayStoreFound {
			fmt.Printf("URL: %s\n", androidapp.PlayStoreInfo.URL)
			fmt.Printf("Version: %s\n", androidapp.PlayStoreInfo.Version)
			fmt.Printf("App summary: %s\n", androidapp.PlayStoreInfo.Summary)
			fmt.Printf("App developer: %s (email: %s, id: %s)\n", androidapp.PlayStoreInfo.Developer, androidapp.PlayStoreInfo.DeveloperEmail, androidapp.PlayStoreInfo.DeveloperID)
			fmt.Printf("Release date: %s\n", androidapp.PlayStoreInfo.Released)
			fmt.Printf("Installs: %s\n", androidapp.PlayStoreInfo.Installs)
			fmt.Printf("Score: %f\n", androidapp.PlayStoreInfo.Score)
		} else {
			fmt.Println("app not found in Play Store")
		}
	}

	fmt.Println()
}
