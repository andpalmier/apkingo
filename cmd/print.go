package main

import (
	"fmt"
	"github.com/fatih/color"
	"reflect"
)

// colors to improve readability
var italic = color.New(color.FgWhite, color.Italic)
var yellow = color.New(color.FgYellow)
var cyan = color.New(color.FgCyan)

// printStruct() - print items in a struct with appropriate
// spacing and colors
func printStruct(s interface{}) {
	items := reflect.ValueOf(s)
	typeOfS := items.Type()

	for i := 0; i < items.NumField(); i++ {
		if len(typeOfS.Field(i).Name) > 6 {
			fmt.Printf("%s:\t", typeOfS.Field(i).Name)
		} else {
			fmt.Printf("%s:\t\t", typeOfS.Field(i).Name)
		}
		if items.Field(i).Interface() == "" {
			italic.Printf("not found\n")
		} else {
			cyan.Printf("%v\n", items.Field(i).Interface())
		}
	}
}

// printGeneralInfo() - print general info section of the apk
func (androidapp *AndroidApp) printGeneralInfo() {
	yellow.Println("\n* General info")
	printStruct(androidapp.GeneralInfo)
}

// printHash() - print hashes of the apk
func (androidapp *AndroidApp) printHash() {
	yellow.Println("\n* Hash values")
	printStruct(androidapp.Hashes)
}

// printPlayStoreInfo() - print Play Store info of the apk
func (androidapp *AndroidApp) printPlayStoreInfo() {
	yellow.Println("\n* Play Store")
	if androidapp.PlayStore.Url != "" {
		printStruct(androidapp.Hashes)
	} else {
		italic.Println("app not found in Play Store")
	}
}

// printCertInfo() - print certificate info of the apk
func (androidapp *AndroidApp) printCertInfo() {
	yellow.Println("\n* Certificate")
	if androidapp.Certificate.Issuer != "" {
		printStruct(androidapp.Certificate)
	} else {
		italic.Println("certificate not found")
	}
}

// printKoodousInfo() - print certificate info of the apk
func (androidapp *AndroidApp) printKoodousInfo() {
	yellow.Println("\n* Koodous")
	if androidapp.Koodous.Url != "" {
		fmt.Printf("Analyzed:\t")
		cyan.Printf("%v\n", androidapp.Koodous.Analyzed)
		fmt.Printf("Detected:\t")
		if androidapp.Koodous.Detected {
			color.Red("[!] apk detected as malicious\n")
		} else {
			cyan.Printf("false\n")
		}
		fmt.Printf("URL:\t\t")
		cyan.Printf("%s\n", androidapp.Koodous.Url)
	} else {
		italic.Printf("impossible to retrieve koodous info\n")
	}
}

// printPermissions() - print permissions found in the apk
func (androidapp *AndroidApp) printPermissions() {
	yellow.Println("\n* Permissions")
	if len(androidapp.Permissions) == 0 {
		italic.Println("no permissions found")
	} else {
		for _, permission := range androidapp.Permissions {
			cyan.Printf("%s\n", permission)
		}
	}
}

// printMetadata() - print metadata found in the apk
func (androidapp *AndroidApp) printMetadata() {
	yellow.Println("\n* Metadata")
	if len(androidapp.Metadata) == 0 {
		italic.Println("no metadata found")
	} else {
		for _, m := range androidapp.Metadata {
			metaname := m.MetadataName
			metavalue := m.MetadataValue
			if metaname != "" {
				if metavalue != "" {
					fmt.Printf("%s: ", metaname)
					cyan.Printf("%s\n", metavalue)
				} else {
					fmt.Printf("%s\n", metaname)
				}
			}
		}
	}
}

// printAndroidInfo() - print all the information
func (androidapp *AndroidApp) printAll() {
	name := androidapp.Name
	yellow.Printf("\nApp name:\t")
	if name != "" {
		cyan.Printf("%s\n", name)
	} else {
		italic.Printf("app name not found\n")
	}
	androidapp.printGeneralInfo()
	androidapp.printHash()
	androidapp.printPermissions()
	androidapp.printMetadata()
	androidapp.printCertInfo()
	androidapp.printPlayStoreInfo()
	androidapp.printKoodousInfo()
	fmt.Println()
}
