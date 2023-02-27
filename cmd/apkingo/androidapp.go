package main

import (
	"fmt"

	"github.com/shogo82148/androidbinary/apk"
)

// permissionsInfo(apk) - get the permission from apk
func permissionsInfo(apk apk.Apk) {
	if len(apk.Manifest().UsesPermissions) == 0 {
		italic.Println("no permissions found")
	} else {
		for _, n := range apk.Manifest().UsesPermissions {
			permission, _ := n.Name.String()
			if permission != "" {
				fmt.Println(permission)
			}
		}
	}
}

// metadataInfo(apk) - get the metadata from apk
func metadataInfo(apk apk.Apk) {
	if len(apk.Manifest().App.MetaData) == 0 {
		italic.Println("no metadata found")
	} else {
		for _, n := range apk.Manifest().App.MetaData {
			metaname, _ := n.Name.String()
			metavalue, _ := n.Value.String()
			fmt.Printf("%s: ", metaname)
			if metavalue != "" {
				cyan.Printf("%s", metavalue)
			}
			fmt.Printf("\n")
		}
	}
}

// getGeneralInfo(apk) - get general info from apk
func generalInfo(apk apk.Apk) {

	yellow.Printf("\nApp name:\t")
	name, err := apk.Label(nil)
	if err == nil {
		cyan.Printf("%s\n", name)
	} else {
		italic.Printf("app name not found\n")
	}

	yellow.Println("\n* General Info")

	fmt.Printf("PackageName:\t")
	printer(apk.PackageName())

	fmt.Printf("App version:\t")
	version, err := apk.Manifest().VersionName.String()
	if err != nil {
		version = ""
	}
	printer(version)

	fmt.Printf("Main activity:\t")
	mainactivity, err := apk.MainActivity()
	if err != nil {
		mainactivity = ""
	}
	printer(mainactivity)

	fmt.Printf("Minimum SDK:\t")
	sdkmin, err := apk.Manifest().SDK.Min.Int32()
	if err != nil {
		italic.Println("not found")
	} else {
		cyan.Printf("%d (%s)\n", sdkmin, androidname[int(sdkmin)])
	}

	fmt.Printf("Target SDK:\t")
	sdktarget, err := apk.Manifest().SDK.Target.Int32()
	if err != nil {
		italic.Println("not found")
	} else {
		cyan.Printf("%d (%s)\n", sdktarget, androidname[int(sdktarget)])
	}
}
