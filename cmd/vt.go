package main

import (
	"errors"
	"fmt"
	"time"
	"encoding/json"
	"github.com/VirusTotal/vt-go"
)

// getVTDetection() - get apk info from VirusTotal using sha256 hash and store them in the androidapp struct
func (androidapp *AndroidApp) getVTDetection(apiKey string) error {
	client := vt.NewClient(apiKey)
	vtitem := VirusTotalInfo{}
	file, err := client.GetObject(vt.URL("files/" + androidapp.Hashes.Sha256))
	if err != nil {
		return errors.New("error performing VirusTotal request")
	} else {
		vtitem.Url = fmt.Sprintf("https://virustotal.com/gui/file/%s", androidapp.Hashes.Sha256)
	}

	fsd, err := file.GetTime("first_submission_date")
	if err == nil {
		vtitem.FirstSubmit = fsd.Format(time.RFC822Z)
	}

	lad, err := file.GetTime("last_analysis_date")
	if err == nil {
		vtitem.LastAnalysis = lad.Format(time.RFC822Z)
	}

	nsubmit, err := file.GetInt64("times_submitted")
	if err == nil {
		vtitem.TimesSubmit = nsubmit
	}

	reputation, err := file.GetInt64("reputation")
	if err == nil {
		vtitem.Reput = reputation
	}
	
	votes, err := file.Get("total_votes")
	if err == nil {
		if votesmap := votes.(map[string]interface{}); votesmap != nil {
			vtitem.Votes.Harmless, _ = votesmap["harmless"].(json.Number).Int64()
			vtitem.Votes.Malicious, _ = votesmap["malicious"].(json.Number).Int64()
		}
	}

	las, err := file.Get("last_analysis_stats")
	if err == nil {
		if lasmap := las.(map[string]interface{}); lasmap != nil {
			vtitem.AnalysStats.Harmless, _ = lasmap["harmless"].(json.Number).Int64()
			vtitem.AnalysStats.Malicious, _ = lasmap["malicious"].(json.Number).Int64()
			vtitem.AnalysStats.TypeUnsupported, _ = lasmap["type-unsupported"].(json.Number).Int64()
			vtitem.AnalysStats.Suspicious, _ = lasmap["suspicious"].(json.Number).Int64()
			vtitem.AnalysStats.ConfirmedTimeout, _ = lasmap["confirmed-timeout"].(json.Number).Int64()
			vtitem.AnalysStats.Timeout, _ = lasmap["timeout"].(json.Number).Int64()
			vtitem.AnalysStats.Failure, _ = lasmap["failure"].(json.Number).Int64()
			vtitem.AnalysStats.Undetected, _ = lasmap["undetected"].(json.Number).Int64()
		}
	}

	names, err := file.GetStringSlice("names")
	if err == nil {
		vtitem.Names = names
	}

	icon, err := file.Get("main_icon")
	if err == nil {
		if iconmap := icon.(map[string]interface{}); iconmap != nil {
			vtitem.Icon.Md5 = fmt.Sprintf("%s",iconmap["raw_md5"])
			vtitem.Icon.Dhash  = fmt.Sprintf("%s",iconmap["dhash"])
		}
	}

	androidapp.VirusTotal = &vtitem
	return nil

}