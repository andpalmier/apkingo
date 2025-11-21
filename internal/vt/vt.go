// Package vt provides integration with the VirusTotal API.
package vt

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/VirusTotal/vt-go"
	"github.com/andpalmier/apkingo/internal/utils"
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
	Package            string      `json:"package"`
	AndroidVersionCode string      `json:"android-version-code"`
	AndroidVersionName string      `json:"android-version-name"`
	MinSdkVersion      string      `json:"min-sdk-version"`
	TargetSdkVersion   string      `json:"target-sdk-version"`
	MainActivity       string      `json:"main-activity"`
	Activities         []string    `json:"activities"`
	Services           []string    `json:"services"`
	Providers          []string    `json:"providers"`
	Receivers          []string    `json:"receivers"`
	Libraries          []string    `json:"libraries"`
	Certificate        interface{} `json:"certificate"`
	DangerPerm         []string    `json:"dangerous-permissions"`
	StringsInformation interface{} `json:"strings-information"`
}

// VirusTotalInfo - struct for info gathered from VirusTotal
type VirusTotalInfo struct {
	Url                      string          `json:"virustotal-url"`
	Names                    []string        `json:"names"`
	SubmissDate              time.Time       `json:"first-submitted"`
	TimesSubmit              int64           `json:"times-submitted"`
	LastAnalysis             time.Time       `json:"last-analysis"`
	AnalysStats              *VTAnalysStats  `json:"analysis-stats"`
	Votes                    *VTVotes        `json:"votes"`
	Icon                     *VTIcon         `json:"icon"`
	Androguard               *AndroguardInfo `json:"androguard"`
	Tags                     []string        `json:"tags"`
	PopularThreatCategory    string          `json:"popular-threat-category"`
	PopularThreatName        string          `json:"popular-threat-name"`
	Reputation               int64           `json:"reputation"`
	TotalCrowdsourcedSigma   int64           `json:"crowdsourced-sigma-analysis-results"`
	TotalCrowdsourcedYara    int64           `json:"crowdsourced-yara-results"`
	TotalCrowdsourcedIDSHits int64           `json:"crowdsourced-ids-results"`
}

// ScanFile uploads and scans a file on VirusTotal
func ScanFile(path string, vtapikey string) error {
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
		utils.LogError("error opening file", err)
		return err
	}
	defer f.Close()

	analysis, err := s.Scan(f, utils.GetFileName(path), progressCh)
	if err != nil {
		utils.LogError("error scanning file", err)
		return err
	}

	fmt.Println("\n\nFile uploaded. The analysis may take some time.")
	fmt.Printf("\nTrack the analysis progress at: https://www.virustotal.com/gui/file-analysis/%s\n\n", analysis.ID())
	return nil
}

// GetInfo retrieves information from VirusTotal using VT API and sha256 hash
func GetInfo(apiKey, hash string) (*VirusTotalInfo, error) {
	client := vt.NewClient(apiKey)
	//nolint:govet // vt.URL is not a printf-style function, string concatenation is intentional
	file, err := client.GetObject(vt.URL("files/%s", hash))
	if err != nil {
		utils.LogError("error getting object from VirusTotal", err)
		return nil, err
	}

	vtinfo := &VirusTotalInfo{
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

	// Additional fields
	if tags, err := file.GetStringSlice("tags"); err == nil {
		vtinfo.Tags = tags
	}

	if popularThreatCategory, err := file.Get("popular_threat_classification"); err == nil {
		if threatMap, ok := popularThreatCategory.(map[string]interface{}); ok {
			if suggestedThreatLabel, ok := threatMap["suggested_threat_label"].(string); ok {
				vtinfo.PopularThreatCategory = suggestedThreatLabel
			}
			if popularThreatName, ok := threatMap["popular_threat_name"].([]interface{}); ok {
				if len(popularThreatName) > 0 {
					if name, ok := popularThreatName[0].(map[string]interface{}); ok {
						if value, ok := name["value"].(string); ok {
							vtinfo.PopularThreatName = value
						}
					}
				}
			}
		}
	}

	if reputation, err := file.GetInt64("reputation"); err == nil {
		vtinfo.Reputation = reputation
	}

	// Crowdsourced detections
	if sigmaResults, err := file.GetInt64("crowdsourced_sigma_analysis_results"); err == nil {
		vtinfo.TotalCrowdsourcedSigma = sigmaResults
	}

	if yaraResults, err := file.GetInt64("crowdsourced_yara_results"); err == nil {
		vtinfo.TotalCrowdsourcedYara = yaraResults
	}

	if idsResults, err := file.GetInt64("crowdsourced_ids_results"); err == nil {
		vtinfo.TotalCrowdsourcedIDSHits = idsResults
	}

	return vtinfo, nil
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
		Package:            getString(data, "Package"),
		AndroidVersionCode: getString(data, "AndroidVersionCode"),
		AndroidVersionName: getString(data, "AndroidVersionName"),
		MinSdkVersion:      getString(data, "MinSdkVersion"),
		TargetSdkVersion:   getString(data, "TargetSdkVersion"),
		MainActivity:       getString(data, "main_activity"),
		Activities:         getStringSlice(data, "Activities"),
		Services:           getStringSlice(data, "Services"),
		Providers:          getStringSlice(data, "Providers"),
		Receivers:          getStringSlice(data, "Receivers"),
		Libraries:          getStringSlice(data, "Libraries"),
		Certificate:        data["certificate"],
		StringsInformation: data["StringsInformation"],
		DangerPerm:         make([]string, 0),
	}

	if permDetails, ok := data["permission_details"].(map[string]interface{}); ok {
		for perm, details := range permDetails {
			if detailsMap, ok := details.(map[string]interface{}); ok {
				if detailsMap["permission_type"] == "dangerous" {
					andro.DangerPerm = append(andro.DangerPerm, perm)
				}
			}
		}
	}

	return &andro
}

func getStringSlice(data map[string]interface{}, key string) []string {
	if val, ok := data[key]; ok {
		if slice, ok := val.([]interface{}); ok {
			result := make([]string, 0, len(slice))
			for _, item := range slice {
				if str, ok := item.(string); ok {
					result = append(result, str)
				}
			}
			return result
		}
	}
	return []string{}
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
