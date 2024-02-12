package main

import (
	"encoding/json"
	"fmt"
	"github.com/VirusTotal/vt-go"
	"os"
	"time"
)

// VTAnalysStats - struct for analysis details by VirusTotal
type VTAnalysStats struct {
	Harmless         int64 `json:"harmless"`
	TypeUnsupported  int64 `json:"type-unsupported"`
	Suspicious       int64 `json:"suspicious"`
	ConfirmedTimeout int64 `json:"confirmed-timeout"`
	Timeout          int64 `json:"timeout"`
	Failure          int64 `json:"failure"`
	Malicious        int64 `json:"malicious"`
	Undetected       int64 `json:"undetected"`
}

// VTVotes - struct for vote details by VirusTotal
type VTVotes struct {
	Harmless  int64 `json:"harmless"`
	Malicious int64 `json:"malicious"`
}

// VTIcon - struct for icon details by VirusTotal
type VTIcon struct {
	Md5   string `json:"md5"`
	Dhash string `json:"dhash"`
}

// AndroguardInfo - struct for info gathered from Androguard
type AndroguardInfo struct {
	Providers  interface{} `json:"providers"`
	Receivers  interface{} `json:"receivers"`
	Services   interface{} `json:"services"`
	Strings    interface{} `json:"strings"`
	DangerPerm []string    `json:"dangerous-permissions"`
}

// VirusTotalInfo - struct for info gathered from VirusTotal
type VirusTotalInfo struct {
	Url          string          `json:"virustotal-url"`
	Names        []string        `json:"names"`
	SubmissDate  time.Time       `json:"first-submitted"`
	TimesSubmit  int64           `json:"times-submitted"`
	LastAnalysis time.Time       `json:"last-analysis"`
	AnalysStats  *VTAnalysStats  `json:"analysis-stats"`
	Votes        *VTVotes        `json:"votes"`
	Icon         *VTIcon         `json:"icon"`
	Androguard   *AndroguardInfo `json:"androguard"`
}

// vtScanFile uploads and scans a file on VirusTotal
func vtScanFile(path string, vtapikey string) error {
	progressCh := make(chan float32)
	defer close(progressCh)

	// Print uploading and scanning progress
	go func() {
		for progress := range progressCh {
			status := "uploading..."
			if progress >= 100 {
				status = "scanning..."
			}
			fmt.Printf("\r%s %s %4.1f%%", path, status, progress)
		}
	}()

	client := vt.NewClient(vtapikey)
	s := client.NewFileScanner()
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	analysis, err := s.Scan(f, getFileName(path), progressCh)
	if err != nil {
		return err
	}

	fmt.Println("\n\nFile uploaded. The analysis may take some time.")
	fmt.Printf("\nTrack the analysis progress at: https://www.virustotal.com/gui/file-analysis/%s\n\n", analysis.ID())
	return nil
}

// setVTInfo retrieves information from VirusTotal using VT API and sha256 hash
func (androidapp *AndroidApp) setVTInfo(apiKey string) error {
	hash := androidapp.Hashes.Sha1
	client := vt.NewClient(apiKey)
	file, err := client.GetObject(vt.URL("files/" + hash))
	if err != nil {
		return err
	}

	vtinfo := VirusTotalInfo{
		Url: fmt.Sprintf("https://virustotal.com/gui/file/%s", hash),
	}

	if names, err := file.GetStringSlice("names"); err == nil {
		vtinfo.Names = names
	}

	if submissDate, err := file.GetTime("first_submission_date"); err == nil {
		vtinfo.SubmissDate = submissDate
	}

	if timesSubmit, err := file.GetInt64("times_submitted"); err == nil {
		vtinfo.TimesSubmit = timesSubmit
	}

	if lastAnalysis, err := file.GetTime("last_analysis_date"); err == nil {
		vtinfo.LastAnalysis = lastAnalysis
	}

	if las, err := file.Get("last_analysis_stats"); err == nil {
		if lasmap, ok := las.(map[string]interface{}); ok {
			vtinfo.AnalysStats = parseVTAnalysStats(lasmap)
		}
	}

	if votes, err := file.Get("total_votes"); err == nil {
		if votesmap, ok := votes.(map[string]interface{}); ok {
			vtinfo.Votes = parseVTVotes(votesmap)
		}
	}

	if icon, err := file.Get("main_icon"); err == nil {
		if iconmap, ok := icon.(map[string]interface{}); ok {
			vtinfo.Icon = parseVTIcon(iconmap)
		}
	}

	if androguard, err := file.Get("androguard"); err == nil {
		if androguardmap, ok := androguard.(map[string]interface{}); ok {
			vtinfo.Androguard = parseAndroguardInfo(androguardmap)
		}
	}

	androidapp.VirusTotal = &vtinfo
	return nil
}

func parseVTAnalysStats(data map[string]interface{}) *VTAnalysStats {
	return &VTAnalysStats{
		Harmless:         getInt64(data, "harmless"),
		TypeUnsupported:  getInt64(data, "type-unsupported"),
		Suspicious:       getInt64(data, "suspicious"),
		ConfirmedTimeout: getInt64(data, "confirmed-timeout"),
		Timeout:          getInt64(data, "timeout"),
		Failure:          getInt64(data, "failure"),
		Malicious:        getInt64(data, "malicious"),
		Undetected:       getInt64(data, "undetected"),
	}
}

func parseVTVotes(data map[string]interface{}) *VTVotes {
	return &VTVotes{
		Harmless:  getInt64(data, "harmless"),
		Malicious: getInt64(data, "malicious"),
	}
}

func parseVTIcon(data map[string]interface{}) *VTIcon {
	return &VTIcon{
		Md5:   getString(data, "raw_md5"),
		Dhash: getString(data, "dhash"),
	}
}

func parseAndroguardInfo(data map[string]interface{}) *AndroguardInfo {
	andro := AndroguardInfo{
		Providers:  data["Providers"],
		Receivers:  data["Receivers"],
		Services:   data["Services"],
		Strings:    data["StringsInformation"],
		DangerPerm: make([]string, 0),
	}

	if permDetails, ok := data["permission_details"].(map[string]interface{}); ok {
		for perm, details := range permDetails {
			if details.(map[string]interface{})["permission_type"] == "dangerous" {
				andro.DangerPerm = append(andro.DangerPerm, perm)
			}
		}
	}

	return &andro
}

func getInt64(data map[string]interface{}, key string) int64 {
	if val, ok := data[key].(json.Number); ok {
		if i, err := val.Int64(); err == nil {
			return i
		}
	}
	return 0
}

func getString(data map[string]interface{}, key string) string {
	if val, ok := data[key].(string); ok {
		return val
	}
	return ""
}
