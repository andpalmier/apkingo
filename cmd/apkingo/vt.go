package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/VirusTotal/vt-go"
)

// printvtnumber(string, number) - help printing numbers retrieved from VT
func printvtnumber(s string, d int64) {
	fmt.Printf("\t%s:\t", s)
	if s == "Malicious" && d > 0 {
		red.Printf("%d - [!] apk detected as malicious\n", d)
	} else {
		cyan.Printf("%d\n", d)
	}
}

// printvtinterface(string, interface) - help printing interfaces retrieved from VT
func printvtinterface(s string, i interface{}) {
	fmt.Printf("%s:\t", s)
	if i != nil {
		cyan.Printf("%v\n", i.([]interface{}))
	} else {
		italic.Printf("not found\n")
	}
}

// vtScanFile(path, vtapikey) - upload and scan file on VT
func vtScanFile(path string, vtapi string) error {
	progressCh := make(chan float32)
	defer close(progressCh)

	// print uploading and scanning progress
	go func() {
		for progress := range progressCh {
			if progress < 100 {
				fmt.Printf("\r%s uploading... %4.1f%%", path, progress)
			} else {
				fmt.Printf("\r%s scanning...", path)
			}
		}
	}()

	client := vt.NewClient(vtapi)
	s := client.NewFileScanner()
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	analysis, err := s.ScanFile(f, progressCh)
	if err != nil {
		return err
	}

	fmt.Println("\n\nFile uploaded, please note that the analysis can take some time.\n" +
		"You can then run apkingo again and see the results of the analysis.")
	fmt.Printf("\nIf you want, you can check the progress of the analysis at: "+
		"https://www.virustotal.com/gui/file-analysis/%s\n\n", analysis.ID())
	return nil
}

// vtInfo(vtAPIkey, hash) - get apk info from VirusTotal using VT API and sha256 hash
func vtInfo(apiKey string, hash string) error {
	client := vt.NewClient(apiKey)
	file, err := client.GetObject(vt.URL("files/" + hash))
	if err != nil {
		return err
	}

	fmt.Printf("URL:\t\thttps://virustotal.com/gui/file/%s\n", hash)

	fmt.Printf("Names:\t\t")
	names, err := file.GetStringSlice("names")
	if err != nil || len(names) == 0 {
		italic.Printf("no names found\n")
	} else {
		cyan.Printf("%v\n", names)
	}

	fsd, _ := file.GetTime("first_submission_date")
	fmt.Printf("Submiss. date:\t")
	cyan.Printf("%s\n", fsd.Format(time.RFC822Z))

	nsubmit, _ := file.GetInt64("times_submitted")
	fmt.Printf("# submissions:\t")
	cyan.Printf("%d\n", nsubmit)

	lad, _ := file.GetTime("last_analysis_date")
	fmt.Printf("Last analysis:\t")
	cyan.Printf("%s\n", lad.Format(time.RFC822Z))

	fmt.Printf("Reputation:\t")
	reputation, err := file.GetInt64("reputation")
	if err == nil {
		italic.Printf("not found\n")
	} else if reputation >= 0 {
		cyan.Printf("%d\n (not malicious)", reputation)
	} else {
		red.Printf("%d\n - [!] apk detected as malicious", reputation)
	}

	fmt.Printf("Last analysis results:\t")
	las, err := file.Get("last_analysis_stats")
	if err == nil {
		if lasmap := las.(map[string]interface{}); lasmap != nil {
			fmt.Println()
			harmless, _ := lasmap["harmless"].(json.Number).Int64()
			printvtnumber("Harmless", harmless)
			malicious, _ := lasmap["malicious"].(json.Number).Int64()
			printvtnumber("Malicious", malicious)
			unsupport, _ := lasmap["type-unsupported"].(json.Number).Int64()
			printvtnumber("Unsupp. type", unsupport)
			suspicious, _ := lasmap["suspicious"].(json.Number).Int64()
			printvtnumber("Suspicious", suspicious)
			confirmedtimeout, _ := lasmap["confirmed-timeout"].(json.Number).Int64()
			printvtnumber("Conf. timeout", confirmedtimeout)
			timeout, _ := lasmap["timeout"].(json.Number).Int64()
			printvtnumber("Timeout", timeout)
			failure, _ := lasmap["failure"].(json.Number).Int64()
			printvtnumber("Failure", failure)
			undetected, _ := lasmap["undetected"].(json.Number).Int64()
			printvtnumber("Undetected", undetected)
		}
	} else {
		italic.Println("not found")
	}

	fmt.Printf("Votes result:\t")
	votes, err := file.Get("total_votes")
	if err == nil {
		if votesmap := votes.(map[string]interface{}); votesmap != nil {
			fmt.Println()
			harmless, _ := votesmap["harmless"].(json.Number).Int64()
			printvtnumber("Harmless", harmless)
			malicious, _ := votesmap["malicious"].(json.Number).Int64()
			printvtnumber("Malicious", malicious)
		}
	} else {
		italic.Println("not found")
	}

	fmt.Printf("Icon hashes:\t")
	icon, err := file.Get("main_icon")
	if err == nil {
		if iconmap := icon.(map[string]interface{}); iconmap != nil {
			fmt.Printf("\n\tmd5:\t\t")
			cyan.Printf("%s\n", iconmap["raw_md5"])
			fmt.Printf("\tdhash:\t\t")
			cyan.Printf("%s\n", iconmap["dhash"])
		}
	} else {
		italic.Println("\ticon not found")
	}

	fmt.Printf("Androguard:\t")
	androguard, err := file.Get("androguard")
	if err == nil {
		if androguardmap := androguard.(map[string]interface{}); androguard != nil {
			fmt.Println()
			printvtinterface("\tProviders", androguardmap["Providers"])
			printvtinterface("\tReceivers", androguardmap["Receivers"])
			printvtinterface("\tServices", androguardmap["Services"])
			printvtinterface("\tStrings", androguardmap["StringsInformation"])

			if androguardmap["permission_details"] != nil {
				fmt.Printf("\tDang. permiss.:\t")
				var dangerpermissions []string
				for i, v := range androguardmap["permission_details"].(map[string]interface{}) {
					if v.(map[string]interface{})["permission_type"] == "dangerous" {
						dangerpermissions = append(dangerpermissions, i)
					}
				}
				if len(dangerpermissions) > 0 {
					red.Printf("%v\n", dangerpermissions)
				} else {
					italic.Printf("not found")
				}
			}
		}
	} else {
		italic.Println("Androguard details not found")
	}

	return nil
}
