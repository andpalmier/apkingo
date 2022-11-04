package main

import (
	"fmt"
	"reflect"

	"github.com/fatih/color"
)

// colors to improve readability
var italic = color.New(color.FgWhite, color.Italic)
var yellow = color.New(color.FgYellow)
var cyan = color.New(color.FgCyan)
var red = color.New(color.FgRed)

// printStruct() - print items in a struct with appropriate
// spacing and colors
func printStruct(s interface{}, prefix string) {
	items := reflect.ValueOf(s)
	types := items.Type()
	for i := 0; i < items.NumField(); i++ {
		name := types.Field(i).Name
		field := items.Field(i).Interface()
		if len(name) > 6 {
			fmt.Printf("%s%s:\t", prefix, name)
		} else {
			fmt.Printf("%s%s:\t\t", prefix, name)
		}

		leng := -1
		fieldarr, check := field.([]interface{})
		if check {
			leng = len(fieldarr)
		}

		if field == "" || leng == 0 {
			italic.Printf("not found\n")
		} else if name == "Detected" && field == true {
			red.Printf("true\n")
		} else if name == "Rating" && field.(int64) < 0 {
			red.Printf("%v\n", field)
		} else if name == "Malicious" && field.(int64) > 0 {
			red.Printf("%v\n", field)
		} else if name == "DangerPermis" {
			red.Printf("%v\n", field)
		} else {
			cyan.Printf("%v\n", field)
		}
	}
}

// printGeneralInfo() - print general info section of the apk
func (androidapp *AndroidApp) printGeneralInfo() {
	yellow.Println("\n* General info")
	printStruct(androidapp.GeneralInfo, "")
}

// printHash() - print hashes
func (androidapp *AndroidApp) printHash() {
	yellow.Println("\n* Hash values")
	printStruct(androidapp.Hashes, "")
}

// printPlayStoreInfo() - print Play Store info
func (androidapp *AndroidApp) printPlayStoreInfo() {
	yellow.Println("\n* Play Store")
	if androidapp.PlayStore != nil {
		printStruct(*androidapp.PlayStore, "")
	} else {
		italic.Println("app not found in Play Store")
	}
}

// printCertInfo() - print certificate info
func (androidapp *AndroidApp) printCertInfo() {
	yellow.Println("\n* Certificate")
	if androidapp.Certificate.Issuer != "" {
		printStruct(androidapp.Certificate, "")
	} else {
		italic.Println("certificate not found")
	}
}

// printPermissions() - print permissions found in the apk
func (androidapp *AndroidApp) printPermissions() {
	yellow.Println("\n* Permissions")
	if len(androidapp.Permissions) == 0 {
		italic.Println("no permissions found")
	} else {
		for _, permission := range androidapp.Permissions {
			fmt.Printf("%s\n", permission)
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

// printKoodousInfo() - print analysis results from Koodous
func (androidapp *AndroidApp) printKoodousInfo() {
	yellow.Println("\n* Koodous")
	printStruct(*androidapp.Koodous, "")
}

// printVTInfo() - print VirusTotal info
func (androidapp *AndroidApp) printVTInfo() {
	yellow.Println("\n* VirusTotal")
	if androidapp.VirusTotal.Url != "" {
		fmt.Printf("URL:\t\t")
		cyan.Printf("%s\n", androidapp.VirusTotal.Url)

		fmt.Printf("Names:\t\t")
		if len(androidapp.VirusTotal.Names) > 0 {
			cyan.Printf("%v\n", androidapp.VirusTotal.Names)
		} else {
			italic.Println("not found")
		}

		fmt.Printf("SubmissionDate:\t")
		cyan.Printf("%s\n", androidapp.VirusTotal.FirstSubmit)
		fmt.Printf("#Submissions:\t")
		cyan.Printf("%d\n", androidapp.VirusTotal.TimesSubmit)
		fmt.Printf("LastAnalysis:\t")
		cyan.Printf("%s\n", androidapp.VirusTotal.LastAnalysis)

		fmt.Printf("Reputation:\t")
		reput := androidapp.VirusTotal.Reput
		if reput < 0 {
			color.Red("%d (malicious)\n", reput)
		} else {
			cyan.Printf("%d (not malicious)\n", reput)
		}

		fmt.Printf("Last analysis results\n")
		printStruct(androidapp.VirusTotal.AnalysStats, "\t")
		fmt.Printf("Votes\n")
		printStruct(androidapp.VirusTotal.Votes, "\t")
		fmt.Printf("Icon\n")
		printStruct(androidapp.VirusTotal.Icon, "\t")
		fmt.Printf("Androguard\n")
		printStruct(androidapp.VirusTotal.Androguard, "\t")

	} else {
		italic.Println("app not found in VirusTotal")
	}
}

// printBanner() - like the cool kids
func printBanner() {
	var Banner string = `
	┌─┐┌─┐┬┌─┬┌┐┌┌─┐┌─┐
	├─┤├─┘├┴┐│││││ ┬│ │
	┴ ┴┴  ┴ ┴┴┘└┘└─┘└─┘
by @andpalmier
`
	cyan.Println(Banner)
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
	if kapi != "" {
		androidapp.printKoodousInfo()
	}
	if vtapi != "" {
		androidapp.printVTInfo()
	}
	fmt.Println()
}
