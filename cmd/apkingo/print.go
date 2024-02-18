package main

import (
	"fmt"
	"reflect"

	"github.com/fatih/color"
)

var (
	italic = color.New(color.FgWhite, color.Italic)
	yellow = color.New(color.FgYellow)
	cyan   = color.New(color.FgCyan)
	red    = color.New(color.FgRed)
)

func printBanner() {
	banner := `
	┌─┐┌─┐┬┌─┬┌┐┌┌─┐┌─┐
	├─┤├─┘├┴┐│││││ ┬│ │
	┴ ┴┴  ┴ ┴┴┘└┘└─┘└─┘
by @andpalmier
`
	cyan.Println(banner)
}

func printString(name, value string) {
	fmt.Printf("%s: ", name)
	if value != "" && value != "<nil>" && value != "[]" {
		cyan.Println(value)
	} else {
		italic.Println("not found")
	}
}

func printStruct(s interface{}) {
	items := reflect.ValueOf(s)
	typeOfS := items.Type()

	for i := 0; i < items.NumField(); i++ {
		fmt.Printf("\t%s: ", typeOfS.Field(i).Name)
		value := items.Field(i).Interface()
		if valueStr, ok := value.(string); ok && valueStr == "" {
			italic.Println("not found")
		} else {
			cyan.Printf("%v\n", value)
		}
	}
}

func printGeneralInfo(androidapp *AndroidApp) {
	yellow.Println("* General info")
	printString("name", androidapp.Name)
	printString("package name", androidapp.PackageName)
	printString("version", androidapp.Version)
	printString("main activity", androidapp.MainActivity)
	fmt.Printf("minimum SDK: ")
	if androidapp.MinimumSDK != 0 {
		cyan.Printf("%d (%s)\n", androidapp.MinimumSDK, androidName[int(androidapp.MinimumSDK)])
	} else {
		italic.Println("not found")
	}

	fmt.Printf("target SDK: ")
	if androidapp.TargetSDK != 0 {
		cyan.Printf("%d (%s)\n", androidapp.TargetSDK, androidName[int(androidapp.TargetSDK)])
	} else {
		italic.Println("not found")
	}
}

func printHash(hashes Hashes) {
	yellow.Println("\n* Hash values")
	printString("md5", hashes.Md5)
	printString("sha1", hashes.Sha1)
	printString("sha256", hashes.Sha256)
}

func printPlayStoreInfo(psinfo *PlayStoreInfo) {
	yellow.Println("\n* Play Store")
	if psinfo != nil {
		printString("URL", psinfo.Url)
		printString("version", psinfo.Version)
		printString("released", psinfo.Release)
		printString("updated", psinfo.Updated.Format("Jan 2, 2006"))
		printString("genre", psinfo.Genre)
		printString("summary", psinfo.Summary)
		printString("installs", psinfo.Installs)
		printString("score", fmt.Sprintf("%v", psinfo.Score))
		fmt.Println("developer:")
		printStruct(psinfo.Developer)
	} else {
		italic.Println("app not found in Play Store")
	}
}

func printCertInfo(certinfo CertificateInfo) {
	yellow.Println("\n* Certificate")
	printString("serial", certinfo.Serial)
	printString("thumbprint", certinfo.Thumbprint)
	printString("valid from", certinfo.ValidFrom)
	printString("valid to", certinfo.ValidTo)

	printIssuerOrSubject("issuer", certinfo.Issuer)
	printIssuerOrSubject("subject", certinfo.Subject)
}

func printIssuerOrSubject(title string, certName CertName) {
	fmt.Printf("%s:\n", title)
	printStruct(certName)
}

func printKoodousInfo(kinfo *KoodousInfo) {
	yellow.Println("\n* Koodous")
	if kinfo != nil {
		printString("url", kinfo.Url)
		printString("id", kinfo.Id)
		printString("submission date", kinfo.SubmissionDate)
		printString("icon link", kinfo.IconLink)
		printString("size (bytes)", fmt.Sprintf("%d", kinfo.Size))
		printString("tags", fmt.Sprintf("%v", kinfo.Tags))
		fmt.Printf("detected: ")
		if kinfo.Detected {
			red.Println("true - app detected as malicious")
		} else {
			cyan.Println("false")
		}

		fmt.Printf("rating: ")
		if kinfo.Rating >= 0 {
			cyan.Printf("%d\n", kinfo.Rating)
		} else {
			red.Printf("%d - negative rating\n", kinfo.Rating)
		}

		printString("corrupted", fmt.Sprintf("%t", kinfo.Corrupted))
	} else {
		italic.Println("impossible to retrieve Koodous info")
	}
}

func printPermissions(permissions []string) {
	yellow.Println("\n* Permissions")
	if len(permissions) == 0 {
		italic.Println("no permissions found")
	} else {
		for _, permission := range permissions {
			fmt.Println(permission)
		}
	}
}

func printMetadata(metadata []Metadata) {
	yellow.Println("\n* Metadata")
	if len(metadata) == 0 {
		italic.Println("no metadata found")
	} else {
		for _, m := range metadata {
			if m.Value != "" {
				fmt.Printf("%s: ", m.Name)
				cyan.Printf("%s\n", m.Value)
			} else {
				fmt.Println(m.Name)
			}
		}
	}
}

func printVTLastAnalysisResults(analysStats *VTAnalysStats) {
	fmt.Println("last analysis results:")
	if analysStats != nil {
		printString("\tharmless", fmt.Sprintf("%d", analysStats.Harmless))
		printString("\ttype unsupported", fmt.Sprintf("%d", analysStats.TypeUnsupported))

		fmt.Printf("\tsuspicious: ")
		if susp := analysStats.Suspicious; susp > 0 {
			red.Printf("%d - app flagged as suspicious\n", susp)
		} else {
			cyan.Println("0")
		}

		printString("\tconfirmed timeout", fmt.Sprintf("%d", analysStats.ConfirmedTimeout))
		printString("\tfailure", fmt.Sprintf("%d", analysStats.Failure))

		fmt.Printf("\tmalicious: ")
		if mal := analysStats.Malicious; mal > 0 {
			red.Printf("%d - app flagged as malicious\n", mal)
		} else {
			cyan.Println("0")
		}

		printString("\tundetected", fmt.Sprintf("%d", analysStats.Undetected))
	} else {
		italic.Println("not found")
	}
}

func printVTIcon(icon *VTIcon) {
	fmt.Printf("icon: ")
	if icon != nil {
		fmt.Printf("\n")
		printStruct(*icon)
	} else {
		italic.Println("no icon found")
	}
}

func printVTVotes(votes *VTVotes) {
	fmt.Println("votes:")
	if votes != nil {
		printString("\tharmless", fmt.Sprintf("%d", votes.Harmless))

		fmt.Printf("\tmalicious: ")
		if mal := votes.Malicious; mal > 0 {
			red.Printf("%d - app voted as malicious\n", mal)
		} else {
			cyan.Println("0")
		}
	} else {
		italic.Println("not found")
	}
}

func printVTAndroguard(androguard *AndroguardInfo) {
	fmt.Printf("androguard: ")
	if androguard != nil {
		fmt.Println()
		printString("\tproviders", fmt.Sprintf("%v", androguard.Providers))
		printString("\treceivers", fmt.Sprintf("%v", androguard.Receivers))
		printString("\tservices", fmt.Sprintf("%v", androguard.Services))
		printString("\tstrings", fmt.Sprintf("%v", androguard.Strings))
		if len(androguard.DangerPerm) > 0 {
			fmt.Printf("\tdangerous permissions: ")
			red.Printf("%v\n", androguard.DangerPerm)
		} else {
			italic.Println("dangerous permissions not found")
		}
	} else {
		italic.Println("not found")
	}
}

func printVTInfo(vtinfo *VirusTotalInfo) {
	yellow.Println("\n* VirusTotal")
	if vtinfo != nil {
		printString("url", vtinfo.Url)
		printString("names", fmt.Sprintf("%v", vtinfo.Names))
		printString("submission date", vtinfo.SubmissDate.String())
		printString("submissions", fmt.Sprintf("%d", vtinfo.TimesSubmit))
		printString("last analysis", vtinfo.LastAnalysis.String())

		printVTLastAnalysisResults(vtinfo.AnalysStats)
		printVTVotes(vtinfo.Votes)
		printVTIcon(vtinfo.Icon)
		printVTAndroguard(vtinfo.Androguard)
	} else {
		italic.Println("app not found in VirusTotal")
	}
}

func printAll(androidapp *AndroidApp) {
	printGeneralInfo(androidapp)
	printHash(androidapp.Hashes)
	printPermissions(androidapp.Permissions)
	printMetadata(androidapp.Metadata)
	printCertInfo(androidapp.Certificate)
	printPlayStoreInfo(androidapp.PlayStore)
	printKoodousInfo(androidapp.Koodous)
	printVTInfo(androidapp.VirusTotal)
}
